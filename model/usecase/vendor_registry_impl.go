package usecase

import (
	"github.com/cloud-mesh/object-storage/model"
	"github.com/cloud-mesh/object-storage/model/repository/storage"
)

type storageRegistryImpl struct {
	storages map[string]storage.Repo
}

func NewRegistry() *storageRegistryImpl {
	return &storageRegistryImpl{
		storages: make(map[string]storage.Repo),
	}
}

func (r *storageRegistryImpl) Register(vendor string, storage storage.Repo) {
	r.storages[vendor] = storage
}

func (r *storageRegistryImpl) Get(vendor string) (storage.Repo, error) {
	client, ok := r.storages[vendor]
	if !ok {
		return nil, model.ErrVendorNotRegistered
	}

	return client, nil
}
