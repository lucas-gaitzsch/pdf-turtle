package pdf

import (
	"context"
	"encoding/json"
	"fmt"
	"io"

	"github.com/lucas-gaitzsch/pdf-turtle/config"
	"github.com/lucas-gaitzsch/pdf-turtle/loopback"
	"github.com/lucas-gaitzsch/pdf-turtle/models"
	"github.com/lucas-gaitzsch/pdf-turtle/services"
	"github.com/lucas-gaitzsch/pdf-turtle/services/assetsprovider"
	"github.com/lucas-gaitzsch/pdf-turtle/services/bundles"
	"github.com/lucas-gaitzsch/pdf-turtle/services/htmlparser"
	"github.com/lucas-gaitzsch/pdf-turtle/utils"
	"github.com/lucas-gaitzsch/pdf-turtle/utils/logging"

	"github.com/lucas-gaitzsch/pdf-turtle/services/templating"
	"github.com/rs/zerolog/log"
)

type PdfServiceAbstraction interface {
	PdfFromHtml(data *models.RenderData) (io.Reader, error)
	PdfFromHtmlTemplate(templateData *models.RenderTemplateData) (io.Reader, error)
	PdfFromBundle(bundle *bundles.Bundle, jsonModel string, templateEngine string) (io.Reader, error)
}

type PdfService struct {
	ctx                   context.Context
	rendererService       services.RendererBackgroundService
	assetsProviderService services.AssetsProviderService
	bundleProviderService services.BundleProviderService
	templateService       templating.TemplateServiceAbstraction
	htmlParser            htmlparser.HtmlParser
}

func NewPdfService(requestctx context.Context) PdfServiceAbstraction {
	return &PdfService{
		ctx:                   requestctx,
		rendererService:       getRendererService(requestctx),
		assetsProviderService: getAssetsProviderService(requestctx),
		bundleProviderService: getBundleProviderService(requestctx),
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

	data, err := logging.LogExecutionTimeWithResults("exec template", ps.ctx, func() (*models.RenderData, error) {
		return ps.templateService.ExecuteTemplate(templateData)
	})

	if err != nil {
		panic(err)
	}

	return ps.renderPdf(data)
}

func (ps *PdfService) PdfFromBundle(bundle *bundles.Bundle, jsonModel string, templateEngine string) (io.Reader, error) {
	conf := config.Get(ps.ctx)

	id, cleanup := ps.bundleProviderService.Provide(bundle)
	defer cleanup()

	opt := bundle.GetOptions()
	opt.BasePath = fmt.Sprintf("http://127.0.0.1:%d%s/%s/", conf.LoopbackPort, loopback.BundlePath, id)

	var pdfData io.Reader
	var errRender error

	hasModel := jsonModel != ""
	hasModelLoggingPreparation := log.Debug().Bool("hasModel", hasModel)

	if hasModel {
		hasModelLoggingPreparation.Msg("got model in form data -> render with template engine")

		if templateEngine != "" {
			log.Debug().
				Str("templateEngine", templateEngine).
				Msg("got templateEngine in form data")
		}

		templateData := &models.RenderTemplateData{
			HtmlTemplate:       bundle.GetBodyHtml(),
			HeaderHtmlTemplate: bundle.GetHeaderHtml(),
			FooterHtmlTemplate: bundle.GetFooterHtml(),
			TemplateEngine:     templateEngine,
			RenderOptions:      opt,
		}

		json.Unmarshal([]byte(jsonModel), &templateData.Model)

		pdfData, errRender = ps.PdfFromHtmlTemplate(templateData)
	} else {
		hasModelLoggingPreparation.Msg("no model given with key 'model' in form data -> render plain html")

		data := &models.RenderData{
			Html:          bundle.GetBodyHtml(),
			HeaderHtml:    bundle.GetHeaderHtml(),
			FooterHtml:    bundle.GetFooterHtml(),
			RenderOptions: opt,
		}

		pdfData, errRender = ps.PdfFromHtml(data)
	}

	return pdfData, errRender
}

func (ps *PdfService) renderPdf(data *models.RenderData) (io.Reader, error) {
	ps.preProcessHtmlData(data)

	data.SetDefaults()

	logging.LogExecutionTime("add styles", ps.ctx, func() {
		ps.addDefaultStyleToHeaderAndFooter(data)

		if !data.RenderOptions.ExcludeBuiltinStyles {
			data.Html = utils.AppendStyleToHtml(data.Html, ps.assetsProviderService.GetMergedCss())
		}
	})

	return logging.LogExecutionTimeWithResults("render pdf", ps.ctx, func() (io.Reader, error) {
		return ps.rendererService.RenderAndReceive(*models.NewJob(ps.ctx, data))
	})
}

func (ps *PdfService) preProcessHtmlData(data *models.RenderData) {
	if data.Html == nil {
		return
	}

	if !data.HasHeaderOrFooterHtml() || !data.RenderOptions.ExcludeBuiltinStyles {

		logging.LogExecutionTime("parse dom", ps.ctx, func() {
			ps.htmlParser.Parse(data.Html)
		})

		if !data.HasHeaderOrFooterHtml() {
			// parse header and footer from main html
			logging.LogExecutionTime("pop header and footer from html", ps.ctx, func() {
				headerHtml, footerHtml := ps.htmlParser.PopHeaderAndFooter()
				data.SetHeaderHtml(headerHtml)
				data.SetFooterHtml(footerHtml)
			})
		}

		body, err := logging.LogExecutionTimeWithResults("parse dom", ps.ctx, func() (*string, error) {
			return ps.htmlParser.GetHtml()
		})

		if err == nil {
			data.Html = body
		} else {
			log.Ctx(ps.ctx).Warn().Err(err).Msg("cant get html from parsed dom")
		}
	}
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

func getBundleProviderService(ctx context.Context) services.BundleProviderService {
	return ctx.Value(config.ContextKeyBundleProviderService).(services.BundleProviderService)
}
