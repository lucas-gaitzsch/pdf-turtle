package bundles

import (
	"archive/zip"
	"encoding/json"
	"errors"
	"io"
	"mime/multipart"
	"strings"

	"github.com/lucas-gaitzsch/pdf-turtle/models"
)

const (
	BundleIndexFile   = "index.html"
	BundleHeaderFile  = "header.html"
	BundleFooterFile  = "footer.html"
	BundleOptionsFile = "options.json"
)

type MultipartFileOpener interface {
	Open() (multipart.File, error)
}

type Opener interface {
	Open() (io.ReadCloser, error)
}

type OpenerFileProxy struct {
	MultipartFileOpener MultipartFileOpener
}

func (o *OpenerFileProxy) Open() (io.ReadCloser, error) {
	r, err := o.MultipartFileOpener.Open()
	return r, err
}

type BundleReader interface {
	GetFileByPath(path string) (io.ReadCloser, error)
	GetFileAsStringByPath(path string) (*string, error)
	GetBodyHtml() *string
	GetHeaderHtml() string
	GetFooterHtml() string
	GetOptions() models.RenderOptions
}

type Bundle struct {
	files map[string]Opener
}

// Read files from zip to intern map (path to file).
// This method can be called multiple times to assemble multiple zip bundles to one bundle.
func (b *Bundle) ReadFromZip(file io.ReaderAt, size int64) error {
	z, err := zip.NewReader(file, size)

	if err != nil {
		return err
	}

	for _, f := range z.File {
		b.AddFile(f.Name, f)
	}

	return nil
}

func (b *Bundle) AddFile(path string, file Opener) {
	if b.files == nil {
		b.files = make(map[string]Opener)
	}

	if !strings.Contains(path, "/") &&
		path != BundleIndexFile &&
		path != BundleHeaderFile &&
		path != BundleFooterFile &&
		path != BundleOptionsFile {
		path = "assets/" + path
	}

	b.files[path] = file
}

func (b *Bundle) TestIndexFile() error {
	if _, hasIndexFile := b.files[BundleIndexFile]; !hasIndexFile {
		return errors.New("no index.html file was found on root of bundle")
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
	s, _ := b.GetFileAsStringByPath(BundleIndexFile)
	return s
}

func (b *Bundle) GetHeaderHtml() string {
	s, _ := b.GetFileAsStringByPath(BundleHeaderFile)
	if s == nil {
		return ""
	}
	return *s
}

func (b *Bundle) GetFooterHtml() string {
	s, _ := b.GetFileAsStringByPath(BundleFooterFile)
	if s == nil {
		return ""
	}
	return *s
}

func (b *Bundle) GetOptions() models.RenderOptions {
	opt := models.RenderOptions{}

	f, err := b.GetFileByPath(BundleOptionsFile)
	if err != nil {
		return opt
	}
	defer f.Close()

	json.NewDecoder(f).Decode(&opt)

	opt.IsBundle = true

	return opt
}
