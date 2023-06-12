package headlesschromium

import (
	"bytes"
	"context"
	"errors"
	"io"
	"time"

	"github.com/lucas-gaitzsch/pdf-turtle/config"

	"github.com/rs/zerolog/log"

	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"
)

//TODO: strip html with:  <script\b[^>]*>([\s\S]*?)<\/script>

func NewChromiumBrowser(ctx context.Context) (context.Context, context.CancelFunc) {

	opts := chromedp.DefaultExecAllocatorOptions[:]

	opts = append(
		opts,
		chromedp.Headless,
		chromedp.Flag("headless", true),
		chromedp.Flag("hide-scrollbars", true),
		chromedp.Flag("mute-audio", true),
		chromedp.Flag("no-sandbox", config.Get(ctx).NoSandbox),
	)

	allocCtx, cancelAllocCtx := chromedp.NewExecAllocator(ctx, opts...)

	cctx, cancelCctx := chromedp.NewContext(allocCtx)

	// Keep chromium browser process running
	if err := chromedp.Run(cctx); err != nil {
		log.Error().Err(err).Msg("chromium browser could not be initialized")
		panic(err)
	}

	log.Info().Msg("new browser initialized")

	cancel := func() {
		log.Info().Msg("cancel called: close chromium")
		cancelCctx()
		cancelAllocCtx()
	}

	return cctx, cancel
}

type OptionsConfigureFunc func(params *page.PrintToPDFParams) *page.PrintToPDFParams

func RenderHtmlAsPdf(chromiumAllocCtx context.Context, outerCtx context.Context, location string, html *string, optionsConfigureFunc OptionsConfigureFunc) (io.Reader, error) {
	if html == nil && location == "" {
		return nil, errors.New("html is nil")
	}

	if location == "" {
		location = "about:blank"
	}

	cctx, cancel := chromedp.NewContext(chromiumAllocCtx)
	defer cancel()

	go func() {
		select {
		case <-cctx.Done():
		case <-outerCtx.Done():
			log.Info().Msg("cancel chromium pdf rendering by outer context")
			cancel()
		}
	}()

	var pdfStream io.Reader
	tasks := chromedp.Tasks{
		chromedp.Navigate(location),

		chromedp.ActionFunc(func(cctx context.Context) error {
			lctx, cancelLctx := context.WithCancel(cctx)
			defer cancelLctx()

			done := make(chan bool, 1)

			chromedp.ListenTarget(lctx, func(ev any) {
				if _, ok := ev.(*page.EventLoadEventFired); ok {
					cancelLctx()
					done <- true
				}
			})

			frameTree, err := page.GetFrameTree().Do(cctx)
			if err != nil {
				return err
			}

			if html != nil {
				if err := page.SetDocumentContent(frameTree.Frame.ID, *html).Do(cctx); err != nil {
					return err
				}
			}

			select {
			case <-done:
				return nil
			case <-time.After(time.Duration(config.Get(outerCtx).RenderTimeoutInSeconds) * time.Second):
				return errors.New("render timeout")
			case <-outerCtx.Done():
				return errors.New("canceled by outer ctx")
			}
		}),

		// injectCss(preloadedMergedCss),

		chromedp.ActionFunc(func(cctx context.Context) error {
			params := page.PrintToPDF()
			b, _, err := optionsConfigureFunc(params).Do(cctx)

			pdfStream = bytes.NewReader(b)

			return err
		}),
	}

	if err := chromedp.Run(cctx, runWithTimeOut(outerCtx, tasks)); err != nil {
		return nil, err
	}

	return pdfStream, nil
}

func runWithTimeOut(outerCtx context.Context, tasks chromedp.Tasks) chromedp.ActionFunc {
	timeout := time.Duration(config.Get(outerCtx).RenderTimeoutInSeconds) * time.Second
	return func(ctx context.Context) error {
		timeoutCtx, cancel := context.WithTimeout(ctx, timeout)
		defer cancel()
		return tasks.Do(timeoutCtx)
	}
}
