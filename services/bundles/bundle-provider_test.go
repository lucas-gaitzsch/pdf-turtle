package bundles

import (
	"testing"

	"github.com/google/uuid"
)

func TestProvide(t *testing.T) {
	bps := NewBundleProviderService()
	id, cleanup := bps.Provide(Bundle{})

	if _, ok := bps.bundles[id]; !ok {
		t.Fatal("no bundle with returned id was provided in map")
	}

	cleanup()

	if _, ok := bps.bundles[id]; ok {
		t.Fatal("bundle should be cleaned up but already exists in map")
	}
}

func TestRemove(t *testing.T) {
	bps := NewBundleProviderService()

	id := uuid.New()
	bps.bundles[id] = Bundle{}

	bps.Remove(id)

	if len(bps.bundles) != 0 {
		t.Fatal("bundle should be removed but already exists in map")
	}
}

func TestGetById(t *testing.T) {
	bps := NewBundleProviderService()

	b := Bundle{
		files: make(map[string]Opener),
	}

	id := uuid.New()
	bps.bundles[id] = b

	bGot, ok := bps.GetById(id)

	if !ok || bGot.files == nil {
		t.Fatal("bundle should be removed but already exists in map")
	}
}
