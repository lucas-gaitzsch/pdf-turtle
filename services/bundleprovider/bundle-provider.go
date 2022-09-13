package bundleprovider

import (
	"sync"

	"github.com/google/uuid"
)

type CleanupFunc = func()

func NewBundleProviderService() *BundleProviderService {
	bps := &BundleProviderService{
		bundles: make(map[uuid.UUID]Bundle),
	}
	return bps
}

type BundleProviderService struct {
	bundles map[uuid.UUID]Bundle
	lock    sync.RWMutex
}

func (bps *BundleProviderService) Provide(bundle Bundle) (id uuid.UUID, cleanup CleanupFunc) {
	bps.lock.Lock()
	defer bps.lock.Unlock()

	id = uuid.New()
	bps.bundles[id] = bundle

	cleanup = func() {
		bps.Remove(id)
	}

	return
}

func (bps *BundleProviderService) Remove(id uuid.UUID) {
	bps.lock.Lock()
	defer bps.lock.Unlock()

	delete(bps.bundles, id)
}

func (bps *BundleProviderService) GetById(id uuid.UUID) (Bundle, bool) {
	bps.lock.RLock()
	defer bps.lock.RUnlock()

	b, ok := bps.bundles[id]

	return b, ok
}
