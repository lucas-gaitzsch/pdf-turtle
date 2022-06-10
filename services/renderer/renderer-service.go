package renderer

import (
	"context"
	"errors"
	"io"
	"os"
	"path"
	"time"

	"github.com/rs/zerolog/log"

	"pdf-turtle/config"
	"pdf-turtle/models"
	"pdf-turtle/static-files/embed"
)

type workerSlot struct{}
type workerSlots chan workerSlot

type RendererBackgroundService struct {
	localCtx          context.Context
	localCtxCancel    context.CancelFunc
	workerSlots       workerSlots
	htmlToPdfRenderer HtmlToPdfRendererAbstraction
	Jobs              chan models.Job

	StaticFilesBuiltin  []string
	StaticFilesExternal []string

	preloadedCss []*string

	renderTimeout   time.Duration
	workerInstances int
}

func NewRendererBackgroundService(ctx context.Context) *RendererBackgroundService {
	rbs := new(RendererBackgroundService)

	rbs.workerInstances = config.Get(ctx).WorkerInstances
	rbs.renderTimeout = time.Duration(config.Get(ctx).RenderTimeoutInSeconds) * time.Second

	rbs.StaticFilesBuiltin = []string{
		"default-pdf-styles.css",
		"modern-pdf-style.css",
	}
	rbs.Init(ctx)

	return rbs
}

func (rbs *RendererBackgroundService) Init(outerCtx context.Context) {
	rbs.workerSlots = make(workerSlots, rbs.workerInstances)
	rbs.Jobs = make(chan models.Job)

	rbs.localCtx, rbs.localCtxCancel = context.WithCancel(outerCtx)

	rbs.preloadCssFilesToCache()

	if rbs.htmlToPdfRenderer == nil {
		rbs.htmlToPdfRenderer = NewAsyncHtmlRendererChromium(rbs.localCtx, rbs.preloadedCss)
	}

	go rbs.handleRequests(outerCtx)

	log.
		Info().
		Int("workerInstances", rbs.workerInstances).
		Msgf("render service started with %d", rbs.workerInstances)
}

func (rbs *RendererBackgroundService) preloadCssFilesToCache() {
	preloadedFiles := make([]*string, 0, len(rbs.StaticFilesBuiltin)+len(rbs.StaticFilesExternal))

	for _, file := range rbs.StaticFilesBuiltin {
		b, err := embed.BuiltinFS.ReadFile(file)
		if err != nil {
			log.Warn().Err(err).Msg("could not load file " + file)
			continue
		}

		str := string(b)
		preloadedFiles = append(preloadedFiles, &str)
	}

	for _, file := range rbs.StaticFilesExternal {
		path := path.Join("static-files", "extern", file)
		b, err := os.ReadFile(path)
		if err != nil {
			log.Warn().Err(err).Msg("could not load file " + file)
			continue
		}

		str := string(b)
		preloadedFiles = append(preloadedFiles, &str)
	}

	rbs.preloadedCss = preloadedFiles
}

func (rbs *RendererBackgroundService) acquiredWorker() {
	rbs.workerSlots <- workerSlot{}

	workerCount := len(rbs.workerSlots)

	log.Debug().
		Int("workerCount", workerCount).
		Msgf("renderer worker up: %d", workerCount)
}

func (rbs *RendererBackgroundService) releaseWorker() {
	<-rbs.workerSlots

	workerCount := len(rbs.workerSlots)

	log.Debug().
		Int("workerCount", workerCount).
		Msgf("renderer worker down: %d", workerCount)
}

func (rbs *RendererBackgroundService) handleRequests(outerCtx context.Context) {
	defer rbs.localCtxCancel()
	defer rbs.htmlToPdfRenderer.Close()
	defer func() {
		if r := recover(); r != nil {
			log.
				Warn().
				Interface("recover", r).
				Msg("Recovered in renderer-service; wait 500 ms")

			time.Sleep(500 * time.Millisecond)

			rbs.htmlToPdfRenderer = nil
			rbs.Init(outerCtx)
		}
	}()

	for {
		select {
		case job := <-rbs.Jobs:
			go rbs.doWork(rbs.localCtx, job)

		case <-rbs.localCtx.Done():
			log.Info().Msg("shutting renderer service down")
			return
		}
	}
}

func (rbs *RendererBackgroundService) doWork(ctx context.Context, job models.Job) {
	rbs.acquiredWorker()
	defer rbs.releaseWorker()

	done := make(chan bool, 1)

	go func() {
		defer func() { done <- true }()

		job.RenderData.RenderOptions.SetDefaults()

		res, err := rbs.htmlToPdfRenderer.RenderHtmlAsPdf(job.RequestCtx, job.RenderData)

		if err != nil {
			job.CallbackChan <- nil
			log.Ctx(job.RequestCtx).Error().Err(err).Msg("render service: cant render pdf")
			return
		}

		select {
		case job.CallbackChan <- res:
		case <-time.After(rbs.renderTimeout):
			log.Ctx(job.RequestCtx).Warn().Msg("render service: cancel pdf callback (timeout)")
		}
	}()

	select {
	case <-done:
	case <-ctx.Done():
		log.Ctx(job.RequestCtx).Info().Msg("cancel render task by global context")
	case <-job.RequestCtx.Done():
		log.Ctx(job.RequestCtx).Info().Msg("cancel render task by request context")
	}
}

func (rbs *RendererBackgroundService) RenderAndReceive(job models.Job) (io.Reader, error) {
	rbs.Jobs <- job

	select {
	case pdfBytes := <-job.CallbackChan:
		if pdfBytes != nil {
			return pdfBytes, nil
		} else {
			return nil, errors.New("pdf callback empty")
		}
	case <-time.After(rbs.renderTimeout + 5*time.Second):
		return nil, errors.New("pdf callback timeout")
	}
}

func (rbs *RendererBackgroundService) Close() {
	for i := 0; i < len(rbs.workerSlots); i++ {
		rbs.releaseWorker()
	}

	rbs.htmlToPdfRenderer.Close()
	rbs.localCtxCancel()
}
