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

const magicBodyPaddingInInches = 0.004

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

	r.init()

	return r
}

func (r *HtmlToPdfRendererChromium) init() {
	defer func() {
		if r.LocalCtx.Err() != nil {
			return
		}

		if err := recover(); err != nil {
			log.Warn().Interface("Err", err).Msg("chromium crashed with panic -> new chromium instance")

			r.chromiumCancelFunc()
			r.watcherCtxCancelFunc()

			// re-init
			r.init()
		}
	}()

	r.ChromiumCtx, r.chromiumCancelFunc = headlesschromium.NewChromiumBrowser(r.LocalCtx)

	r.startWatchingChromiumInstance()
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
			log.Warn().Err(r.ChromiumCtx.Err()).Msg("chromium crashed with err -> new chromium instance")

			r.chromiumCancelFunc()
			r.watcherCtxCancelFunc()

			// re-init
			r.init()
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
			WithMarginRight(utils.MmToInches(margins.Right) + magicBodyPaddingInInches).
			WithMarginBottom(utils.MmToInches(margins.Bottom)).
			WithMarginLeft(utils.MmToInches(margins.Left) + magicBodyPaddingInInches)

		if hasHeaderOrFooter {
			var headerFooterWidth int

			if !data.RenderOptions.Landscape {
				headerFooterWidth = data.RenderOptions.PageSize.Width
			} else {
				headerFooterWidth = data.RenderOptions.PageSize.Height
			}

			headerFooterAppendCss := fmt.Sprintf(`
				#header, #footer {
					box-sizing: border-box;
					padding: 0 !important;
					margin:  0 !important;
					
					width: %dmm !important;
					height: %dmm !important;
				}
				#footer {
					height: %dmm !important;				
				}
				#header > div, #footer > div {
					box-sizing: border-box;

					padding-left: %dmm !important;
					padding-right: %dmm !important;

					min-width: %dmm;
					max-width: %dmm;

					min-height: %dmm;
					max-height: %dmm;
				}
				#footer > div {
					transform-origin: bottom left;
					min-height: %dmm;
					max-height: %dmm;
				}
			`,
				headerFooterWidth,
				margins.Top,
				margins.Bottom,
				margins.Left,
				margins.Right,
				headerFooterWidth,
				headerFooterWidth,
				margins.Top,
				margins.Top,
				margins.Bottom,
				margins.Bottom)

			headerHtmlPtr := utils.AppendStyleToHtml(&data.HeaderHtml, &headerFooterAppendCss)
			footerHtmlPtr := utils.AppendStyleToHtml(&data.FooterHtml, &headerFooterAppendCss)

			headerHtmlPtr = utils.RequestAndInlineAllHtmlResources(ctx, headerHtmlPtr, data.RenderOptions.BasePath)
			footerHtmlPtr = utils.RequestAndInlineAllHtmlResources(ctx, footerHtmlPtr, data.RenderOptions.BasePath)

			params = params.
				WithHeaderTemplate(*headerHtmlPtr).
				WithFooterTemplate(*footerHtmlPtr)
		}

		return params
	}

	bodyHtml := data.Html

	return headlesschromium.RenderHtmlAsPdf(r.ChromiumCtx, ctx, data.RenderOptions.BasePath, bodyHtml, paramsFunc)
}

func (r *HtmlToPdfRendererChromium) Close() {
	r.watcherCtxCancelFunc()
	<-r.watcherClosedChan

	r.chromiumCancelFunc()
}
