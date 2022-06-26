package renderer

import (
	"context"
	"fmt"
	"io"

	"github.com/lucas-gaitzsch/pdf-turtle/models"
	"github.com/lucas-gaitzsch/pdf-turtle/services/renderer/headlesschromium"
	"github.com/lucas-gaitzsch/pdf-turtle/utils"

	"github.com/chromedp/cdproto/page"
	"github.com/rs/zerolog/log"
)

// TODO:! Testen was passiert, wenn chromium prozess beendet wurde

const headerFooterScaleFactor = 0.75

type HtmlToPdfRendererChromium struct {
	ChromiumCtx          context.Context
	chromiumCancelFunc   context.CancelFunc
	LocalCtx             context.Context
	watcherCtxCancelFunc context.CancelFunc
	watcherClosedChan    chan bool
}

func NewAsyncHtmlRendererChromium(ctx context.Context) *HtmlToPdfRendererChromium {
	r := new(HtmlToPdfRendererChromium)
	r.LocalCtx = ctx
	r.ChromiumCtx, r.chromiumCancelFunc = headlesschromium.NewChromiumBrowser(r.LocalCtx)

	r.startWatchingChromiumInstance()

	return r
}

func (r *HtmlToPdfRendererChromium) startWatchingChromiumInstance() {
	r.watcherClosedChan = make(chan bool, 1)
	var watcherCtx context.Context
	watcherCtx, r.watcherCtxCancelFunc = context.WithCancel(r.LocalCtx)
	go func() {
		select {
		case <-watcherCtx.Done():
			r.watcherClosedChan <- true
		case <-r.ChromiumCtx.Done():
			log.Warn().Err(r.ChromiumCtx.Err()).Msg("chromium crashed -> new chromium instance")
			r.chromiumCancelFunc()
			r.ChromiumCtx, r.chromiumCancelFunc = headlesschromium.NewChromiumBrowser(r.LocalCtx)

			r.watcherCtxCancelFunc()
			r.startWatchingChromiumInstance()
		}
	}()
}

func (r *HtmlToPdfRendererChromium) RenderHtmlAsPdf(ctx context.Context, data *models.RenderData) (io.Reader, error) {

	hasHeaderOrFooter := data.HasHeaderOrFooterHtml()

	paramsFunc := func(params *page.PrintToPDFParams) *page.PrintToPDFParams {

		margins := models.RenderOptionsMargins{}

		if data.RenderOptions.Margins != nil {
			margins = *data.RenderOptions.Margins
		} else {
			data.RenderOptions.Margins = &margins
		}

		params = params.WithPrintBackground(false).
			WithPreferCSSPageSize(false).
			WithDisplayHeaderFooter(hasHeaderOrFooter).
			WithLandscape(data.RenderOptions.Landscape).
			WithPaperWidth(utils.MmToInches(data.RenderOptions.PageSize.Width)).
			WithPaperHeight(utils.MmToInches(data.RenderOptions.PageSize.Height)).
			WithMarginTop(utils.MmToInches(margins.Top)).
			WithMarginRight(utils.MmToInches(margins.Right)).
			WithMarginBottom(utils.MmToInches(margins.Bottom)).
			WithMarginLeft(utils.MmToInches(margins.Left))

		if hasHeaderOrFooter {
			var headerFooterWidth int

			if !data.RenderOptions.Landscape {
				headerFooterWidth = data.RenderOptions.PageSize.Width
			} else {
				headerFooterWidth = data.RenderOptions.PageSize.Height
			}

			scaledHeaderFooterWidth := float64(headerFooterWidth) * headerFooterScaleFactor

			scaledHeaderHeight := float64(margins.Top) * headerFooterScaleFactor
			scaledFooterHeight := float64(margins.Bottom) * headerFooterScaleFactor

			headerFooterAppendCss := fmt.Sprintf(`
				#header, #footer {
					box-sizing: border-box;
					padding: 0 !important;
					margin:  0 !important;
					
					width: %fmm !important;
					height: %fmm !important;
				}
				#footer {
					height: %fmm !important;				
				}
				#header > div, #footer > div {
					box-sizing: border-box;
					transform: scale(%f);
					transform-origin: top left;

					padding-left: %dmm !important;
					padding-right: %dmm !important;

					min-width: %dmm;
					min-height: %dmm;
				}
				#footer > div {
					transform-origin: bottom left;
					min-height: %dmm;
				}
			`,
				scaledHeaderFooterWidth,
				scaledHeaderHeight,
				scaledFooterHeight,
				headerFooterScaleFactor,
				margins.Left,
				margins.Right,
				headerFooterWidth,
				margins.Top,
				margins.Bottom)

			headerHtml := utils.AppendStyleToHtml(&data.HeaderHtml, &headerFooterAppendCss)
			footerHtml := utils.AppendStyleToHtml(&data.FooterHtml, &headerFooterAppendCss)

			params = params.
				WithHeaderTemplate(*headerHtml).
				WithFooterTemplate(*footerHtml)
		}

		return params
	}

	bodyHtml := data.Html

	return headlesschromium.RenderHtmlAsPdf(r.ChromiumCtx, ctx, bodyHtml, paramsFunc)
}

func (r *HtmlToPdfRendererChromium) Close() {
	r.watcherCtxCancelFunc()
	<-r.watcherClosedChan

	r.chromiumCancelFunc()
}
