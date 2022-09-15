package bundles

import (
	"archive/zip"
	"encoding/json"
	"errors"
	"io"
	"strings"

	"github.com/lucas-gaitzsch/pdf-turtle/models"
)

type Opener interface {
	Open() (io.ReadCloser, error)
}

type Bundle struct {
	files map[string]Opener
}

func (b *Bundle) ReadFromZip(file io.ReaderAt, size int64) error {
	b.files = make(map[string]Opener)

	z, err := zip.NewReader(file, size)

	if err != nil {
		return err
	}

	for _, f := range z.File {
		b.files[f.Name] = f
	}

	return nil
}

func (b *Bundle) GetFileByPath(path string) (io.ReadCloser, error) {
	f, ok := b.files[path]

	if !ok {
		return nil, errors.New("no file found: " + path)
	}

	return f.Open()
}

func (b *Bundle) GetFileAsStringByPath(path string) (*string, error) {
	f, err := b.GetFileByPath(path)

	if err != nil {
		return nil, err
	}
	defer f.Close()

	buf := new(strings.Builder)
	_, err = io.Copy(buf, f)

	if err != nil {
		return nil, err
	}

	str := buf.String()

	return &str, nil
}

func (b *Bundle) GetBodyHtml() *string {
	s, _ := b.GetFileAsStringByPath("index.html")
	return s
}

func (b *Bundle) GetHeaderHtml() string {
	s, _ := b.GetFileAsStringByPath("header.html")
	if s == nil {
		return ""
	}
	return *s
}

func (b *Bundle) GetFooterHtml() string {
	s, _ := b.GetFileAsStringByPath("footer.html")
	if s == nil {
		return ""
	}
	return *s
}

func (b *Bundle) GetOptions() models.RenderOptions {
	f, err := b.GetFileByPath("options.json")

	opt := models.RenderOptions{}

	if err != nil {
		return opt
	}
	defer f.Close()

	json.NewDecoder(f).Decode(&opt)

	opt.IsBundle = true

	return opt
}
