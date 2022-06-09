package renderer

import (
	"context"
	"fmt"
	"io"
	"pdf-turtle/models"
	"pdf-turtle/services/headlesschromium"
	"pdf-turtle/utils"

	"github.com/chromedp/cdproto/page"
	"github.com/rs/zerolog/log"
)

// TODO:! Testen was passiert, wenn chromium prozess beendet wurde

type HtmlToPdfRendererChromium struct {
	ChromiumCtx          context.Context
	chromiumCancelFunc   context.CancelFunc
	LocalCtx             context.Context
	PreloadedMergedCss   *string
	watcherCtxCancelFunc context.CancelFunc
}

func NewAsyncHtmlRendererChromium(ctx context.Context, preloadedCss []*string) *HtmlToPdfRendererChromium {
	r := new(HtmlToPdfRendererChromium)
	r.LocalCtx = ctx
	r.ChromiumCtx, r.chromiumCancelFunc = headlesschromium.NewChromiumBrowser(r.LocalCtx)

	r.PreloadedMergedCss = utils.MergeCss(preloadedCss)

	r.startWatchingChromiumInstance()

	return r
}

func (r *HtmlToPdfRendererChromium) startWatchingChromiumInstance() {
	var watcherCtx context.Context
	watcherCtx, r.watcherCtxCancelFunc = context.WithCancel(r.LocalCtx)
	go func() {
		select {
		case <-watcherCtx.Done():
		case <-r.ChromiumCtx.Done():
			log.Warn().Err(r.ChromiumCtx.Err()).Msg("Chromium crashed -> new Chromium instance")
			r.chromiumCancelFunc()
			r.ChromiumCtx, r.chromiumCancelFunc = headlesschromium.NewChromiumBrowser(r.LocalCtx)

			r.watcherCtxCancelFunc()
			r.startWatchingChromiumInstance()
		}
	}()
}

func (r *HtmlToPdfRendererChromium) RenderHtmlAsPdf(ctx context.Context, data *models.RenderData) (io.Reader, error) {

	hasHeaderOrFooter := data.HeaderHtml != "" || data.FooterHtml != ""

	paramsFunc := func(params *page.PrintToPDFParams) *page.PrintToPDFParams {

		margins := models.RenderOptionsMargins{}

		if data.RenderOptions.Margins != nil {
			margins = *data.RenderOptions.Margins
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
			headerHtml := &data.HeaderHtml
			footerHtml := &data.FooterHtml

			headerFooterAppendCss := fmt.Sprintf(`
				#header, #footer {
					padding: 0 !important;
					padding-left: %dmm !important;
					padding-right: %dmm !important;
					
					transform: scale(0.75);
					transform-origin: top left;
					width: 100%%;
				}
				#footer {
					transform-origin: bottom left;					
				}
			`, margins.Left, margins.Right)

			if !data.RenderOptions.ExcludeBuiltinStyles {
				headerFooterStyle := *r.PreloadedMergedCss + headerFooterAppendCss

				headerHtml = utils.AppendStyleToHtml(headerHtml, &headerFooterStyle)
				footerHtml = utils.AppendStyleToHtml(footerHtml, &headerFooterStyle)
			}

			params = params.
				WithHeaderTemplate(*headerHtml).
				WithFooterTemplate(*footerHtml)
		}

		return params
	}

	bodyHtml := data.BodyHtml

	if !data.RenderOptions.ExcludeBuiltinStyles {
		bodyHtml = utils.AppendStyleToHtml(bodyHtml, r.PreloadedMergedCss)
	}

	return headlesschromium.RenderHtmlAsPdf(r.ChromiumCtx, ctx, bodyHtml, paramsFunc)
}

func (r *HtmlToPdfRendererChromium) Close() {
	r.watcherCtxCancelFunc()
	r.chromiumCancelFunc()
}
