package pdf

import (
	"context"
	"io"

	"github.com/lucas-gaitzsch/pdf-turtle/config"
	"github.com/lucas-gaitzsch/pdf-turtle/models"
	"github.com/lucas-gaitzsch/pdf-turtle/services"
	"github.com/lucas-gaitzsch/pdf-turtle/services/assetsprovider"
	"github.com/lucas-gaitzsch/pdf-turtle/services/htmlparser"
	"github.com/lucas-gaitzsch/pdf-turtle/utils"

	"github.com/lucas-gaitzsch/pdf-turtle/services/templating"
	"github.com/rs/zerolog/log"
)

type PdfServiceAbstraction interface {
	PdfFromHtml(data *models.RenderData) (io.Reader, error)
	PdfFromHtmlTemplate(templateData *models.RenderTemplateData) (io.Reader, error)
}

type PdfService struct {
	ctx                   context.Context
	rendererService       services.RendererBackgroundService
	assetsProviderService services.AssetsProviderService
	templateService       templating.TemplateServiceAbstraction
	htmlParser            htmlparser.HtmlParser
}

func NewPdfService(requestctx context.Context) PdfServiceAbstraction {
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

	data, err := utils.LogExecutionTimeWithResults("exec template", ps.ctx, func() (*models.RenderData, error) {
		return ps.templateService.ExecuteTemplate(templateData)
	})

	if err != nil {
		panic(err)
	}

	return ps.renderPdf(data)
}

func (ps *PdfService) renderPdf(data *models.RenderData) (io.Reader, error) {
	ps.preProcessHtmlData(data)

	data.SetDefaults()

	utils.LogExecutionTime("add styles", ps.ctx, func() {
		ps.addDefaultStyleToHeaderAndFooter(data)

		if !data.RenderOptions.ExcludeBuiltinStyles {
			data.Html = utils.AppendStyleToHtml(data.Html, ps.assetsProviderService.GetMergedCss())
		}
	})

	return utils.LogExecutionTimeWithResults("render pdf", ps.ctx, func() (io.Reader, error) {
		return ps.rendererService.RenderAndReceive(*models.NewJob(ps.ctx, data))
	})
}

func (ps *PdfService) preProcessHtmlData(data *models.RenderData) {
	if data.Html == nil {
		return
	}

	if !data.HasHeaderOrFooterHtml() || !data.RenderOptions.ExcludeBuiltinStyles {

		utils.LogExecutionTime("parse dom", ps.ctx, func() {
			ps.htmlParser.Parse(data.Html)
		})

		if !data.HasHeaderOrFooterHtml() {
			// parse header and footer from main html
			ps.popHeaderAndFooter(data)
		}

		body, err := utils.LogExecutionTimeWithResults("parse dom", ps.ctx, func() (*string, error) {
			return ps.htmlParser.GetHtml()
		})

		if err == nil {
			data.Html = body
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

func getRendererService(ctx context.Context) services.RendererBackgroundService {
	return ctx.Value(config.ContextKeyRendererService).(services.RendererBackgroundService)
}

func getAssetsProviderService(ctx context.Context) services.AssetsProviderService {
	return ctx.Value(config.ContextKeyAssetsProviderService).(services.AssetsProviderService)
}
