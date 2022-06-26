package pdf

import (
	"context"
	"io"

	"github.com/lucas-gaitzsch/pdf-turtle/config"
	"github.com/lucas-gaitzsch/pdf-turtle/models"
	"github.com/lucas-gaitzsch/pdf-turtle/services/assetsprovider"
	"github.com/lucas-gaitzsch/pdf-turtle/services/htmlparser"
	"github.com/lucas-gaitzsch/pdf-turtle/services/renderer"
	"github.com/lucas-gaitzsch/pdf-turtle/utils"

	"github.com/lucas-gaitzsch/pdf-turtle/services/templating"
	"github.com/rs/zerolog/log"
)

type PdfService struct {
	ctx                   context.Context
	rendererService       *renderer.RendererBackgroundService
	assetsProviderService *assetsprovider.AssetsProviderService
	templateService       templating.TemplateServiceAbstraction
	htmlParser            htmlparser.HtmlParser
}

func NewPdfService(requestctx context.Context) *PdfService {
	return &PdfService{
		ctx:                   requestctx,
		rendererService:       getRendererService(requestctx),
		assetsProviderService: getAssetsProviderService(requestctx),
		templateService:       templating.NewTemplateService(),
		htmlParser:            htmlparser.New(),
	}
}

func (ps *PdfService) PdfFromHtml(data *models.RenderData) (io.Reader, error) {
	ps.preProcessHtmlData(data)
	return ps.renderPdf(data)
}

func (ps *PdfService) PdfFromHtmlTemplate(templateData *models.RenderTemplateData) (io.Reader, error) {

	templateData.ParseJsonModelDataFromDoubleEncodedString()

	data, err := utils.LogExecutionTimeWithResult("exec template", ps.ctx, func() (*models.RenderData, error) {
		return ps.templateService.ExecuteTemplate(templateData)
	})

	if err != nil {
		panic(err)
	}

	return ps.renderPdf(data)
}

//TODO:!
// func (ps *PdfService) PdfFromMarkdown() (io.Reader, error) {

// }

func (ps *PdfService) renderPdf(data *models.RenderData) (io.Reader, error) {
	ps.preProcessHtmlData(data)

	data.SetDefaults()

	utils.LogExecutionTime("add styles", ps.ctx, func() {
		ps.addDefaultStyleToHeaderAndFooter(data)

		if !data.HasBuiltinStylesExcluded() {
			htmlWithStyles := utils.AppendStyleToHtml(data.GetBodyHtml(), ps.assetsProviderService.GetMergedCss())
			data.SetBodyHtml(htmlWithStyles)
		}
	})

	return utils.LogExecutionTimeWithResult("render pdf", ps.ctx, func() (io.Reader, error) {
		return ps.rendererService.RenderAndReceive(*models.NewJob(ps.ctx, data))
	})
}

func (ps *PdfService) preProcessHtmlData(data *models.RenderData) {
	if data.GetBodyHtml() == nil {
		return
	}

	if !data.HasHeaderOrFooterHtml() || !data.HasBuiltinStylesExcluded() {

		utils.LogExecutionTime("parse dom", ps.ctx, func() {
			ps.htmlParser.Parse(data.GetBodyHtml())
		})

		if !data.HasHeaderOrFooterHtml() {
			// parse header and footer from main html
			ps.popHeaderAndFooter(data)
		}

		body, err := utils.LogExecutionTimeWithResult("parse dom", ps.ctx, func() (*string, error) {
			return ps.htmlParser.GetHtml()
		})

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
		if headerHtml := data.GetHeaderHtml(); headerHtml != "" {
			data.SetHeaderHtml(*utils.AppendStyleToHtml(&headerHtml, defaultCss))
		}

		if footerHtml := data.GetFooterHtml(); footerHtml != "" {
			data.SetFooterHtml(*utils.AppendStyleToHtml(&footerHtml, defaultCss))
		}
	} else {
		log.Ctx(ps.ctx).Warn().Msg("could not load default styles")
	}
}

func getRendererService(ctx context.Context) *renderer.RendererBackgroundService {
	return ctx.Value(config.ContextKeyRendererService).(*renderer.RendererBackgroundService)
}

func getAssetsProviderService(ctx context.Context) *assetsprovider.AssetsProviderService {
	return ctx.Value(config.ContextKeyAssetsProviderService).(*assetsprovider.AssetsProviderService)
}
