package templating

import (
	"errors"

	"github.com/lucas-gaitzsch/pdf-turtle/models"
	"github.com/lucas-gaitzsch/pdf-turtle/services/templating/templateengines"
)

type TemplateServiceAbstraction interface {
	ExecuteTemplate(data *models.RenderTemplateData) (*models.RenderData, error)
}

func NewTemplateService() TemplateServiceAbstraction {
	return &TemplateService{}
}

type TemplateService struct {
}

func (ts *TemplateService) ExecuteTemplate(templateData *models.RenderTemplateData) (*models.RenderData, error) {
	if templateData == nil {
		return nil, errors.New("template data model should not be nil")
	}

	templateEngine, found := templateengines.GetTemplateEngineByKey(templateData.TemplateEngine)
	
	templateengines.LogParsedTemplateEngine(templateData.TemplateEngine, templateEngine, found)

	data := &models.RenderData{
		RenderOptions: templateData.RenderOptions,
	}

	html, err := templateEngine.Execute(templateData.HtmlTemplate, templateData.Model)
	if err != nil {
		return nil, err
	}

	headerHtml, err := templateEngine.Execute(&templateData.HeaderHtmlTemplate, templateData.Model)
	if err != nil {
		return nil, err
	}

	footerHtml, err := templateEngine.Execute(&templateData.FooterHtmlTemplate, templateData.Model)
	if err != nil {
		return nil, err
	}

	data.Html = html
	data.HeaderHtml = *headerHtml
	data.FooterHtml = *footerHtml

	return data, nil
}
