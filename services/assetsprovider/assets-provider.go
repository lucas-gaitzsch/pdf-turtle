package assetsprovider

import (
	"os"
	"path"
	"pdf-turtle/utils"

	"pdf-turtle/static-files/embed"

	"github.com/rs/zerolog/log"
)

const (
	DefaultPdfStyles = "default-pdf-styles.css"
	PdfTurtleStyles  = "pdf-turtle-styles.css"
)

func NewAssetsProviderService() *AssetsProviderService {
	aps := &AssetsProviderService{}
	aps.init()
	return aps
}

type AssetsProviderService struct {
	StaticFilesBuiltin  []string
	StaticFilesExternal []string

	preloadedCss map[string]*string

	mergedCss *string
}

func (aps *AssetsProviderService) init() {
	aps.preloadedCss = make(map[string]*string)
	aps.StaticFilesBuiltin = []string{
		DefaultPdfStyles,
		PdfTurtleStyles,
	}

	aps.preloadCssFilesToCache()
}

func (aps *AssetsProviderService) preloadCssFilesToCache() {

	for _, file := range aps.StaticFilesBuiltin {
		b, err := embed.BuiltinFS.ReadFile(file)
		if err != nil {
			log.Warn().Err(err).Msg("could not load file " + file)
			continue
		}

		str := string(b)
		aps.preloadedCss[file] = &str
	}

	for _, file := range aps.StaticFilesExternal {
		path := path.Join("static-files", "extern", file)
		b, err := os.ReadFile(path)
		if err != nil {
			log.Warn().Err(err).Msg("could not load file " + file)
			continue
		}

		str := string(b)
		aps.preloadedCss[file] = &str
	}

	preloadedCssArr := make([]*string, 0, len(aps.preloadedCss))
	for _, value := range aps.preloadedCss {
		preloadedCssArr = append(preloadedCssArr, value)
	}

	aps.mergedCss = utils.MergeCss(preloadedCssArr)
}

func (aps *AssetsProviderService) GetMergedCss() *string {
	return aps.mergedCss
}

func (aps *AssetsProviderService) GetCssByKey(key string) (css *string, ok bool) {
	css, ok = aps.preloadedCss[key]
	return
}
