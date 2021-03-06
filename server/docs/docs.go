// Package docs GENERATED BY SWAG; DO NOT EDIT
// This file was generated by swaggo/swag
package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "contact": {
            "name": "Lucas Gaitzsch",
            "email": "lucas@gaitzsch.dev"
        },
        "license": {
            "name": "Apache 2.0",
            "url": "http://www.apache.org/licenses/LICENSE-2.0.html"
        },
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/pdf/from/html-template/render": {
            "post": {
                "description": "Returns PDF file generated from HTML template plus model of body, header and footer",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/pdf"
                ],
                "tags": [
                    "render html-template"
                ],
                "summary": "Render PDF from HTML template",
                "parameters": [
                    {
                        "description": "Render Data",
                        "name": "renderTemplateData",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/RenderTemplateData"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "PDF File"
                    }
                }
            }
        },
        "/pdf/from/html-template/test": {
            "post": {
                "description": "Returns information about matching model data to template",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "test html-template"
                ],
                "summary": "Test HTML template matching model",
                "parameters": [
                    {
                        "description": "Render Data",
                        "name": "renderTemplateData",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/RenderTemplateData"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/TemplateTestResult"
                        }
                    }
                }
            }
        },
        "/pdf/from/html/render": {
            "post": {
                "description": "Returns PDF file generated from HTML of body, header and footer",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/pdf"
                ],
                "tags": [
                    "render html"
                ],
                "summary": "Render PDF from HTML",
                "parameters": [
                    {
                        "description": "Render Data",
                        "name": "renderData",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/RenderData"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "PDF File"
                    }
                }
            }
        }
    },
    "definitions": {
        "PageSize": {
            "type": "object",
            "properties": {
                "height": {
                    "description": "in mm",
                    "type": "integer",
                    "example": 297
                },
                "width": {
                    "description": "in mm",
                    "type": "integer",
                    "example": 210
                }
            }
        },
        "RenderData": {
            "type": "object",
            "properties": {
                "footerHtml": {
                    "description": "Optional html for footer. If empty, the footer html will be parsed from main html (\u003cPdfFooter\u003e\u003c/PdfFooter\u003e).",
                    "type": "string",
                    "default": "\u003cdiv class=\"default-footer\"\u003e\u003cdiv\u003e\u003cspan class=\"pageNumber\"\u003e\u003c/span\u003e of \u003cspan class=\"totalPages\"\u003e\u003c/span\u003e\u003c/div\u003e\u003c/div\u003e"
                },
                "headerHtml": {
                    "description": "Optional html for header. If empty, the header html will be parsed from main html (\u003cPdfHeader\u003e\u003c/PdfHeader\u003e).",
                    "type": "string",
                    "example": "\u003ch1\u003eHeading\u003c/h1\u003e"
                },
                "html": {
                    "type": "string",
                    "example": "\u003cb\u003eHello World\u003c/b\u003e"
                },
                "options": {
                    "$ref": "#/definitions/RenderOptions"
                }
            }
        },
        "RenderOptions": {
            "type": "object",
            "properties": {
                "excludeBuiltinStyles": {
                    "type": "boolean",
                    "default": false
                },
                "landscape": {
                    "type": "boolean",
                    "default": false
                },
                "margins": {
                    "description": "margins in mm; fallback to default if null",
                    "$ref": "#/definitions/RenderOptionsMargins"
                },
                "pageFormat": {
                    "type": "string",
                    "default": "A4",
                    "enum": [
                        "A0",
                        "A1",
                        "A2",
                        "A3",
                        "A4",
                        "A5",
                        "A6",
                        "Letter",
                        "Legal"
                    ]
                },
                "pageSize": {
                    "description": "page size in mm; overrides page format",
                    "$ref": "#/definitions/PageSize"
                }
            }
        },
        "RenderOptionsMargins": {
            "type": "object",
            "properties": {
                "bottom": {
                    "description": "margin bottom in mm",
                    "type": "integer",
                    "default": 20
                },
                "left": {
                    "description": "margin left in mm",
                    "type": "integer",
                    "default": 25
                },
                "right": {
                    "description": "margin right in mm",
                    "type": "integer",
                    "default": 25
                },
                "top": {
                    "description": "margin top in mm",
                    "type": "integer",
                    "default": 25
                }
            }
        },
        "RenderTemplateData": {
            "type": "object",
            "properties": {
                "footerHtmlTemplate": {
                    "description": "Optional template for footer. If empty, the footer template will be parsed from main template (\u003cPdfFooter\u003e\u003c/PdfFooter\u003e).",
                    "type": "string"
                },
                "footerModel": {
                    "description": "Optional model for footer. If empty or null model was used.",
                    "type": "object"
                },
                "headerHtmlTemplate": {
                    "description": "Optional template for header. If empty, the header template will be parsed from main template (\u003cPdfHeader\u003e\u003c/PdfHeader\u003e).",
                    "type": "string"
                },
                "headerModel": {
                    "description": "Optional model for header. If empty or null model was used.",
                    "type": "object"
                },
                "htmlTemplate": {
                    "type": "string"
                },
                "model": {
                    "type": "object"
                },
                "options": {
                    "$ref": "#/definitions/RenderOptions"
                },
                "templateEngine": {
                    "type": "string",
                    "default": "golang",
                    "enum": [
                        "golang",
                        "handlebars",
                        "django"
                    ]
                }
            }
        },
        "TemplateTestResult": {
            "type": "object",
            "properties": {
                "bodyTemplateError": {
                    "type": "string"
                },
                "footerTemplateError": {
                    "type": "string"
                },
                "headerTemplateError": {
                    "type": "string"
                },
                "isValid": {
                    "type": "boolean"
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "",
	BasePath:         "/api",
	Schemes:          []string{"http"},
	Title:            "PdfTurtle API",
	Description:      "A painless HTML to PDF rendering service. Generate PDF reports and documents from HTML templates or raw HTML.",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
