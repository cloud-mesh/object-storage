package usecase

import (
	"github.com/cloud-mesh/object-storage/model"
	log "github.com/sirupsen/logrus"
	"io"
)

type usecaseImpl struct {
	storageRegistry StorageRegistry
}

func NewUseCase(storageRegistry StorageRegistry) *usecaseImpl {
	return &usecaseImpl{storageRegistry}
}

func (u *usecaseImpl) GetBucket(vendor string, bucketName string) (*model.Bucket, error) {
	storage, err := u.storageRegistry.Get(vendor)
	if err != nil {
		return nil, err
	}

	bucket, err := storage.GetBucket(bucketName)
	if err != nil {
		log.WithError(err).Errorf("vendor get bucket: vendor=%s, bucket=%s", vendor, bucketName)
		return nil, err
	}

	return bucket, nil
}

func (u *usecaseImpl) ListBucket(vendor string, page int, pageSize int) ([]*model.Bucket, error) {
	storage, err := u.storageRegistry.Get(vendor)
	if err != nil {
		return nil, err
	}

	buckets, err := storage.ListBucket(page, pageSize)
	if err != nil {
		log.WithError(err).Errorf("vendor list bucket: vendor=%s, page=%d, pageSize=%d", vendor, page, pageSize)
		return nil, err
	}

	return buckets, nil
}

func (u *usecaseImpl) CreateBucket(vendor string, bucketName string) error {
	storage, err := u.storageRegistry.Get(vendor)
	if err != nil {
		return err
	}

	if err := storage.CreateBucket(bucketName); err != nil {
		log.WithError(err).Errorf("vendor create bucket: vendor=%s, bucket=%s", vendor, bucketName)
		return err
	}

	return nil
}

func (u *usecaseImpl) DeleteBucket(vendor string, bucketName string) error {
	storage, err := u.storageRegistry.Get(vendor)
	if err != nil {
		return err
	}

	if err := storage.DeleteBucket(bucketName); err != nil {
		log.WithError(err).Errorf("vendor delete bucket: vendor=%s, bucket=%s", vendor, bucketName)
		return err
	}

	return nil
}

func (u *usecaseImpl) HeadObject(vendor string, bucketName string, objectKey string) (*model.Object, error) {
	storage, err := u.storageRegistry.Get(vendor)
	if err != nil {
		return nil, err
	}

	object, err := storage.HeadObject(bucketName, objectKey)
	if err != nil {
		log.WithError(err).Errorf("vendor head object: vendor=%s, bucket=%s, objectKey=%s", vendor, bucketName, objectKey)
		return nil, err
	}
	return object, nil
}

func (u *usecaseImpl) PutObject(vendor string, bucketName string, objectKey string, reader io.Reader) error {
	storage, err := u.storageRegistry.Get(vendor)
	if err != nil {
		return err
	}

	if err := storage.PutObject(bucketName, objectKey, reader); err != nil {
		log.WithError(err).Errorf("vendor put object: vendor=%s, bucket=%s, objectKey=%s", vendor, bucketName, objectKey)
		return err
	}

	return nil
}

func (u *usecaseImpl) DeleteObject(vendor string, bucketName string, objectKey string) error {
	storage, err := u.storageRegistry.Get(vendor)
	if err != nil {
		return err
	}

	if err := storage.DeleteObject(bucketName, objectKey); err != nil {
		log.WithError(err).Errorf("vendor delete object: vendor=%s, bucket=%s, objectKey=%s", vendor, bucketName, objectKey)
		return err
	}

	return nil
}

func (u *usecaseImpl) InitMultipartUpload(vendor string, bucketName string, objectKey string) (string, error) {
	storage, err := u.storageRegistry.Get(vendor)
	if err != nil {
		return "", err
	}

	uploadID, err := storage.InitMultipartUpload(bucketName, objectKey)
	if err != nil {
		log.WithError(err).Errorf("vendor init multipart: vendor=%s, bucket=%s, objectKey=%s", vendor, bucketName, objectKey)
		return "", nil
	}

	return uploadID, nil
}

func (u *usecaseImpl) CompleteUploadPart(vendor string, bucketName string, objectKey string, uploadId string, parts []model.ObjectCompletePart) error {
	storage, err := u.storageRegistry.Get(vendor)
	if err != nil {
		return err
	}

	if err := storage.CompleteUploadPart(bucketName, objectKey, uploadId, parts); err != nil {
		log.WithError(err).Errorf("vendor complete multipart: vendor=%s, bucket=%s, objectKey=%s, uploadId=%s, parts=%#v",
			vendor, bucketName, objectKey, uploadId, parts)
		return err
	}

	return nil
}

func (u *usecaseImpl) AbortMultipartUpload(vendor string, bucketName string, objectKey string, uploadId string) error {
	storage, err := u.storageRegistry.Get(vendor)
	if err != nil {
		return err
	}

	if err := storage.AbortMultipartUpload(bucketName, objectKey, uploadId); err != nil {
		log.WithError(err).Errorf("vendor abort multipart: vendor=%s, bucket=%s, objectKey=%s, uploadId=%s",
			vendor, bucketName, objectKey, uploadId)
		return err
	}

	return nil
}

func (u *usecaseImpl) UploadPart(vendor string, bucketName string, objectKey string, uploadId string, partNum int, reader io.ReadSeeker) (string, error) {
	storage, err := u.storageRegistry.Get(vendor)
	if err != nil {
		return "", err
	}

	eTag, err := storage.UploadPart(bucketName, objectKey, uploadId, partNum, reader)
	if err != nil {
		log.WithError(err).Errorf("vendor upload multipart: vendor=%s, bucket=%s, objectKey=%s, uploadId=%s, partNumber=%d",
			vendor, bucketName, objectKey, uploadId, partNum)
		return "", err
	}

	return eTag, nil
}

func (u *usecaseImpl) ListParts(vendor string, bucketName string, objectKey string, uploadId string) ([]*model.ObjectPart, error) {
	storage, err := u.storageRegistry.Get(vendor)
	if err != nil {
		return nil, err
	}

	parts, err := storage.ListParts(bucketName, objectKey, uploadId)
	if err != nil {
		log.WithError(err).Errorf("vendor list multipart: vendor=%s, bucket=%s, objectKey=%s, uploadId=%s",
			vendor, bucketName, objectKey, uploadId)
		return nil, err
	}

	return parts, nil
}
