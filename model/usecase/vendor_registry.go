package usecase

import (
	"github.com/cloud-mesh/object-storage/model/repository/storage"
)

type StorageRegistry interface {
	Register(vendor string, storage storage.Repo)
	Get(vendor string) (storage storage.Repo, err error)
}
