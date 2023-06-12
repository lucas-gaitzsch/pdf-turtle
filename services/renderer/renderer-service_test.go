package renderer

import (
	"bytes"
	"context"
	"io"
	"testing"
	"time"

	"github.com/lucas-gaitzsch/pdf-turtle/models"
	"github.com/lucas-gaitzsch/pdf-turtle/utils/logging"
)

type htmlToPdfRendererMock struct {
	ContinueChan  chan bool
	HitRenderChan chan bool
}

func (m *htmlToPdfRendererMock) RenderHtmlAsPdf(ctx context.Context, data *models.RenderData) (io.Reader, error) {
	m.HitRenderChan <- true
	<-m.ContinueChan
	return bytes.NewReader([]byte{}), nil
}

func (m *htmlToPdfRendererMock) Close() {
	close(m.ContinueChan)
	close(m.HitRenderChan)
}

func newTestRenderService(ctx context.Context, workerInstances int) (*RendererBackgroundService, *htmlToPdfRendererMock) {
	rbs := new(RendererBackgroundService)

	rendererMock := &htmlToPdfRendererMock{
		ContinueChan:  make(chan bool),
		HitRenderChan: make(chan bool),
	}

	rbs.htmlToPdfRenderer = rendererMock

	rbs.workerInstances = workerInstances
	rbs.renderTimeout = 5 * time.Second

	rbs.Init(ctx)

	return rbs, rendererMock
}

func TestWorkerUpAndDown(t *testing.T) {
	logging.InitTestLogger(t)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	service, rendererMock := newTestRenderService(ctx, 10)
	defer service.Close()

	gotReturn := make(chan bool)

	go func() {
		_, _ = service.RenderAndReceive(*models.NewJob(context.Background(), &models.RenderData{}))
		gotReturn <- true
	}()

	<-rendererMock.HitRenderChan

	if len(service.workerSlots) != 1 {
		t.Fatal("worker slots should have len of 1 while request")
	}

	rendererMock.ContinueChan <- true

	<-gotReturn

	for i := 0; len(service.workerSlots) > 0 || i < 100; i++ {
		// busy wait for all workers shutdown
		<-time.After(50 * time.Millisecond)
	}

	if len(service.workerSlots) != 0 {
		t.Fatal("worker slots should have len of 0 after return")
	}
}

func TestWorkerBeyondTheLimit(t *testing.T) {
	const jobCount = 50
	const workerInstances = 40

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	service, rendererMock := newTestRenderService(ctx, workerInstances)
	defer service.Close()

	gotReturn := make(chan bool)

	// queue jobs
	for i := 0; i < jobCount; i++ {
		go func() {
			service.RenderAndReceive(*models.NewJob(context.Background(), &models.RenderData{}))
			gotReturn <- true
		}()
	}

	// check worker limit
	for i := 0; i < workerInstances; i++ {
		<-rendererMock.HitRenderChan
	}

	if currLen := len(service.workerSlots); currLen != workerInstances {
		t.Fatalf("worker slots should have len of %d while request (curr: %d)", workerInstances, currLen)
	}

	for i := 0; i < workerInstances; i++ {
		rendererMock.ContinueChan <- true
		<-gotReturn
	}

	// process remaining
	for i := 0; i < jobCount-workerInstances; i++ {
		<-rendererMock.HitRenderChan
	}
	
	<-time.After(50 * time.Millisecond)

	if currLen := len(service.workerSlots); currLen != jobCount-workerInstances {
		t.Fatalf("worker slots should have len of %d while request (curr: %d)", jobCount-workerInstances, currLen)
	}

	for i := 0; i < jobCount-workerInstances; i++ {
		rendererMock.ContinueChan <- true
		<-gotReturn
	}

	for i := 0; len(service.workerSlots) > 0 || i < 100; i++ {
		// busy wait for all workers shutdown
		<-time.After(50 * time.Millisecond)
	}

	if currLen := len(service.workerSlots); currLen != 0 {
		t.Fatalf("worker slots should have len of 0 after return (curr: %d)", currLen)
	}
}
