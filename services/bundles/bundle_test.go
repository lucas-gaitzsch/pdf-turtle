package bundles

import (
	"bytes"
	"os"
	"strings"
	"testing"
)

const assetsTestFile = "assets/testasset.png"

func TestReadFromZip(t *testing.T) {
	b := getTestBundle()

	if len(b.files) != 5 {
		t.Fatal("parsed test-bundle should be contain 5 files (index.html + header.html + footer.html + options.json + testasset.png)")
	}

	assertHasAssetAndCanReadFile(t, b, BundleIndexFile)
	assertHasAssetAndCanReadFile(t, b, BundleHeaderFile)
	assertHasAssetAndCanReadFile(t, b, BundleFooterFile)
	assertHasAssetAndCanReadFile(t, b, BundleOptionsFile)
	assertHasAssetAndCanReadFile(t, b, assetsTestFile)
}

func TestReadFromZipMultiple(t *testing.T) {
	f, _ := os.ReadFile("../../test-assets/test-bundle.zip")

	r := bytes.NewReader(f)

	b := &Bundle{
		files: make(map[string]Opener),
	}
	b.files["test.html"] = nil
	b.files["index.html"] = nil
	b.ReadFromZip(r, int64(len(f)))

	if _, ok := b.files["test.html"]; !ok {
		t.Fatal("existing file was omitted")
	}

	if bf, ok := b.files["index.html"]; !ok || bf == nil {
		t.Fatal("index file should be overwritten")
	}
}

func TestGetFileByPath(t *testing.T) {
	b := getTestBundle()

	r, err := b.GetFileByPath(assetsTestFile)
	if err != nil {
		t.Fatalf("err should be null but is: %v", err)
	}
	r.Close()
}

func TestGetFileAsStringByPath(t *testing.T) {
	b := getTestBundle()

	sPtr, err := b.GetFileAsStringByPath(BundleHeaderFile)
	if err != nil {
		t.Fatalf("err should be null but is: %v", err)
	}
	s := *sPtr
	if !strings.Contains(s, "test header") {
		t.Fatalf("header content should contains 'test header' but is: %s", s)
	}
}

func TestGetBody(t *testing.T) {
	b := getTestBundle()

	sPtr := b.GetBodyHtml()
	s := *sPtr
	if !strings.Contains(s, "test body") {
		t.Fatalf("header content should contains 'test body' but is: %s", s)
	}
}

func TestGetHeader(t *testing.T) {
	b := getTestBundle()

	s := b.GetHeaderHtml()
	if !strings.Contains(s, "test header") {
		t.Fatalf("header content should contains 'test header' but is: %s", s)
	}
}

func TestGetFooter(t *testing.T) {
	b := getTestBundle()

	s := b.GetFooterHtml()
	if !strings.Contains(s, "test footer") {
		t.Fatalf("header content should contains 'test footer' but is: %s", s)
	}
}

func TestGetOptions(t *testing.T) {
	b := getTestBundle()
	opt := b.GetOptions()

	if opt.Landscape != true {
		t.Fatal("options should contains true landscape flag")
	}
}

func getTestBundle() *Bundle {
	f, _ := os.ReadFile("../../test-assets/test-bundle.zip")

	r := bytes.NewReader(f)

	b := &Bundle{}
	b.ReadFromZip(r, int64(len(f)))

	return b
}

func assertHasAssetAndCanReadFile(t *testing.T, b *Bundle, name string) {
	opener, hasFile := b.files[name]
	if !hasFile {
		t.Fatalf("parsed bundle should contain '%s' file", name)
	}

	r, err := opener.Open()
	if err != nil {
		t.Fatalf("cant open file")
	}
	defer r.Close()

	buffer := make([]byte, 5)
	n, err := r.Read(buffer)
	if err != nil || n == 0 {
		t.Fatalf("cant read anything: %s; err: %v", name, err)
	}
}
