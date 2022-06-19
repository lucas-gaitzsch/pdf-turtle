package pdf

import (
	"context"
	"io"
	"pdf-turtle/config"
	"pdf-turtle/models"
	"pdf-turtle/services/assetsprovider"
	"pdf-turtle/services/htmlparser"
	"pdf-turtle/services/renderer"
	"pdf-turtle/services/templating/templateengines"
	"pdf-turtle/utils"

	"github.com/rs/zerolog/log"
)

type PdfService struct {
	ctx                   context.Context
	rendererService       *renderer.RendererBackgroundService
	assetsProviderService *assetsprovider.AssetsProviderService
	htmlParser            htmlparser.HtmlParser
}

func NewPdfService(requestctx context.Context) *PdfService {
	return &PdfService{
		ctx:                   requestctx,
		rendererService:       getRendererService(requestctx),
		assetsProviderService: getAssetsProviderService(requestctx),
		htmlParser:            htmlparser.New(),
	}
}

func (ps *PdfService) PdfFromHtml(data *models.RenderData) (io.Reader, error) {
	ps.preProcessHtmlData(data)
	return ps.renderPdf(data)
}

func (ps *PdfService) PdfFromHtmlTemplate(templateData *models.RenderTemplateData) (io.Reader, error) {

	templateData.ParseJsonModelDataFromDoubleEncodedString()

	templateEngine := templateengines.GetTemplateEngineByKey(templateData.TemplateEngine)

	data := &models.RenderData{
		RenderOptions: templateData.RenderOptions,
	}

	utils.LogExecutionTime("exec template", ps.ctx, func() {
		html, err := templateEngine.Execute(templateData.HtmlTemplate, templateData.Model)
		if err != nil {
			panic(err)
		}

		headerHtml, err := templateEngine.Execute(&templateData.HeaderHtmlTemplate, templateData.GetHeaderModel())
		if err != nil {
			panic(err)
		}

		footerHtml, err := templateEngine.Execute(&templateData.FooterHtmlTemplate, templateData.GetFooterModel())
		if err != nil {
			panic(err)
		}

		data.Html = html
		data.HeaderHtml = *headerHtml
		data.FooterHtml = *footerHtml
	})

	return ps.renderPdf(data)
}

//TODO:!
// func (ps *PdfService) PdfFromMarkdown() (io.Reader, error) {

// }

func (ps *PdfService) preProcessHtmlData(data *models.RenderData) {
	if data.GetBodyHtml() == nil {
		return
	}

	if !data.HasHeaderOrFooterHtml() || !data.HasBuiltinStylesExcluded() {

		ps.htmlParser.Parse(data.GetBodyHtml())

		if !data.HasHeaderOrFooterHtml() {
			// parse header and footer from main html
			ps.popHeaderAndFooter(data)
		}

		ps.addDefaultStyleToHeaderAndFooter(data)

		if !data.HasBuiltinStylesExcluded() {
			ps.htmlParser.AddStyles(ps.assetsProviderService.GetMergedCss())
		}

		body, err := ps.htmlParser.GetHtml()

		if err == nil {
			data.SetBodyHtml(body)
		} else {
			log.Ctx(ps.ctx).Warn().Err(err).Msg("cant get html from parsed dom")
		}
	}
}

func (ps *PdfService) popHeaderAndFooter(data HtmlModels) {
	utils.LogExecutionTime("pop header and footer from html", ps.ctx, func() {
		headerHtml, footerHtml := ps.htmlParser.PopHeaderAndFooter()
		data.SetHeaderHtml(headerHtml)
		data.SetFooterHtml(footerHtml)
	})
}

func (ps *PdfService) addDefaultStyleToHeaderAndFooter(data HtmlModels) {
	defaultCss, ok := ps.assetsProviderService.GetCssByKey(assetsprovider.DefaultPdfStyles)
	if ok {
		headerHtml := data.GetHeaderHtml()
		data.SetHeaderHtml(*utils.AppendStyleToHtml(&headerHtml, defaultCss))

		footerHtml := data.GetFooterHtml()
		data.SetFooterHtml(*utils.AppendStyleToHtml(&footerHtml, defaultCss))
	} else {
		log.Ctx(ps.ctx).Warn().Msg("could not load default styles")
	}
}

func (ps *PdfService) renderPdf(data *models.RenderData) (io.Reader, error) {
	ps.preProcessHtmlData(data)

	data.SetDefaults()

	return utils.LogExecutionTimeWithResult("render pdf", ps.ctx, func() (io.Reader, error) {
		return ps.rendererService.RenderAndReceive(*models.NewJob(ps.ctx, data))
	})
}

func getRendererService(ctx context.Context) *renderer.RendererBackgroundService {
	return ctx.Value(config.ContextKeyRendererService).(*renderer.RendererBackgroundService)
}

func getAssetsProviderService(ctx context.Context) *assetsprovider.AssetsProviderService {
	return ctx.Value(config.ContextKeyAssetsProviderService).(*assetsprovider.AssetsProviderService)
}
