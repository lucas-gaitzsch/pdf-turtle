package config

type ContextKey string

const (
	ContextKeyConfig                = ContextKey("config")
	ContextKeyRendererService       = ContextKey("rendererService")
	ContextKeyAssetsProviderService = ContextKey("assetsProviderService")
	ContextKeyRequestId             = ContextKey("requestId")
)
