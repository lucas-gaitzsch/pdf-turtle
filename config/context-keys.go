package config

type ContextKey string

const (
	ContextKeyConfig     = ContextKey("config")
	ContextKeyPdfService = ContextKey("pdfService")
	ContextKeyRequestId  = ContextKey("requestId")
)
