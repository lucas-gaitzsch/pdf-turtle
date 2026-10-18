package headlesschromium

import (
	"bytes"
	"context"
	"errors"
	"io"
	"time"

	"github.com/lucas-gaitzsch/pdf-turtle/config"

	"github.com/rs/zerolog/log"

	"github.com/chromedp/cdproto/cdp"
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

	if config.Get(ctx).NoSandbox {
		opts = append(
			opts,
			chromedp.Flag("no-sandbox", true),
		)
	}

	allocCtx, cancelAllocCtx := chromedp.NewExecAllocator(ctx, opts...)

	cctx, cancelCctx := chromedp.NewContext(allocCtx)

	// Keep chromium browser process running
	if err := chromedp.Do(cctx); err != nil {
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

type OptionsConfigureFunc func(params *page.PrintToPDFParams)

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
	tasks := []chromedp.Action[chromedp.Void]{
		chromedp.Navigate(location),
	}

	if html != nil {
		tasks = append(tasks, chromedp.Func(func(ctx context.Context, target *chromedp.Target) error {
			_, err := chromedp.WaitEvent(page.LoadEventFired, nil,
				chromedp.Func(func(ctx context.Context, target *chromedp.Target) error {
					frameTree, err := cdp.Call(ctx, target, page.GetFrameTree, cdp.Empty{})
					if err != nil {
						return err
					}

					_, err = cdp.Call(ctx, target, page.SetDocumentContent, page.SetDocumentContentParams{
						FrameID: frameTree.FrameTree.Frame.ID,
						HTML:    *html,
					})
					return err
				}),
			)(ctx, target)
			return err
		}))
	}

	// injectCss(preloadedMergedCss)

	tasks = append(tasks, chromedp.Func(func(ctx context.Context, target *chromedp.Target) error {
		params := page.PrintToPDFParams{}
		optionsConfigureFunc(&params)

		result, err := cdp.Call(ctx, target, page.PrintToPDF, params)
		pdfStream = bytes.NewReader(result.Data)
		return err
	}))

	if err := chromedp.Do(cctx, runWithTimeOut(outerCtx, chromedp.Steps(tasks...))); err != nil {
		return nil, err
	}

	return pdfStream, nil
}

func runWithTimeOut(outerCtx context.Context, tasks chromedp.Action[chromedp.Void]) chromedp.Action[chromedp.Void] {
	timeout := time.Duration(config.Get(outerCtx).RenderTimeoutInSeconds) * time.Second
	return func(ctx context.Context, target *chromedp.Target) (chromedp.Void, error) {
		timeoutCtx, cancel := context.WithTimeout(ctx, timeout)
		defer cancel()
		return tasks(timeoutCtx, target)
	}
}
