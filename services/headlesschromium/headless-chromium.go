package headlesschromium

import (
	"bytes"
	"context"
	"errors"
	"io"
	"pdf-turtle/config"
	"time"

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
		log.Warn().Msg("cancel called: close chromium")
		cancelCctx()
		cancelAllocCtx()
	}

	return cctx, cancel
}

type OptionsConfigureFunc func(params *page.PrintToPDFParams) *page.PrintToPDFParams

func RenderHtmlAsPdf(chromiumAllocCtx context.Context, outerCtx context.Context, html *string, optionsConfigureFunc OptionsConfigureFunc) (io.Reader, error) {
	if html == nil {
		return nil, errors.New("html is nil")
	}

	cctx, cancel := chromedp.NewContext(chromiumAllocCtx)
	defer cancel()

	go func() {
		select {
		case <-cctx.Done():
			log.Debug().Msg("cancel cctx")
		case <-outerCtx.Done():
			log.Info().Msg("cancel chromium pdf rendering by outer context")
			cancel()
		}
	}()

	var pdfStream io.Reader
	tasks := chromedp.Tasks{
		chromedp.Navigate("about:blank"),

		chromedp.ActionFunc(func(cctx context.Context) error {
			lctx, cancelLctx := context.WithCancel(cctx)
			defer cancelLctx()

			done := make(chan bool, 1)

			chromedp.ListenTarget(lctx, func(ev interface{}) {
				if _, ok := ev.(*page.EventLoadEventFired); ok {
					cancelLctx()
					done <- true
				}
			})

			frameTree, err := page.GetFrameTree().Do(cctx)
			if err != nil {
				return err
			}

			if err := page.SetDocumentContent(frameTree.Frame.ID, *html).Do(cctx); err != nil {
				return err
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
		log.Error().Err(err).Msg("pdf render err")
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

func injectCss(css *string) chromedp.Action {
	if css == nil {
		log.Warn().Msg("css to inject is nil")
	}

	// https://github.com/chromedp/chromedp/issues/520#issuecomment-836446035
	// TODO: a bit hacky .. could be better ...

	const script = `
	(css) => {
		const style = document.createElement('style');
		style.type = 'text/css';
		style.appendChild(document.createTextNode(css));
		document.head.appendChild(style);

		return true;
	}
	`

	return chromedp.PollFunction(script, nil, chromedp.WithPollingArgs(css))
}
