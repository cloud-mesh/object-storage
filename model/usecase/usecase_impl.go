package usecase

import (
	"github.com/cloud-mesh/object-storage/model"
	"github.com/cloud-mesh/object-storage/model/repository/storage"
	log "github.com/sirupsen/logrus"
	"io"
)

type usecaseImpl struct {
	vendor  string
	storage storage.Repo
}

func NewUseCase(vendor string, storage storage.Repo) *usecaseImpl {
	return &usecaseImpl{vendor, storage}
}

func (u *usecaseImpl) HeadRegion() (region *model.Region, err error) {
	return &model.Region{Vendor: u.vendor}, nil
}

func (u *usecaseImpl) GetBucket(bucketName string) (*model.Bucket, error) {
	bucket, err := u.storage.GetBucket(bucketName)
	if err != nil {
		log.WithError(err).Errorf("vendor get bucket: bucket=%s", bucketName)
		return nil, err
	}

	return bucket, nil
}

func (u *usecaseImpl) ListBucket(page int, pageSize int) ([]*model.Bucket, error) {
	buckets, err := u.storage.ListBucket(page, pageSize)
	if err != nil {
		log.WithError(err).Errorf("vendor list bucket: page=%d, pageSize=%d", page, pageSize)
		return nil, err
	}

	return buckets, nil
}

func (u *usecaseImpl) CreateBucket(bucketName string) error {
	if err := u.storage.CreateBucket(bucketName); err != nil {
		log.WithError(err).Errorf("vendor create bucket: bucket=%s", bucketName)
		return err
	}

	return nil
}

func (u *usecaseImpl) DeleteBucket(bucketName string) error {
	if err := u.storage.DeleteBucket(bucketName); err != nil {
		log.WithError(err).Errorf("vendor delete bucket: bucket=%s", bucketName)
		return err
	}

	return nil
}

func (u *usecaseImpl) HeadObject(bucketName string, objectKey string) (*model.Object, error) {
	object, err := u.storage.HeadObject(bucketName, objectKey)
	if err != nil {
		log.WithError(err).Errorf("vendor head object: bucket=%s, objectKey=%s", bucketName, objectKey)
		return nil, err
	}
	return object, nil
}

func (u *usecaseImpl) PutObject(bucketName string, objectKey string, reader io.Reader) error {
	if err := u.storage.PutObject(bucketName, objectKey, reader); err != nil {
		log.WithError(err).Errorf("vendor put object: bucket=%s, objectKey=%s", bucketName, objectKey)
		return err
	}

	return nil
}

func (u *usecaseImpl) DeleteObject(bucketName string, objectKey string) error {
	if err := u.storage.DeleteObject(bucketName, objectKey); err != nil {
		log.WithError(err).Errorf("vendor delete object: bucket=%s, objectKey=%s", bucketName, objectKey)
		return err
	}

	return nil
}

func (u *usecaseImpl) InitMultipartUpload(bucketName string, objectKey string) (string, error) {
	uploadID, err := u.storage.InitMultipartUpload(bucketName, objectKey)
	if err != nil {
		log.WithError(err).Errorf("vendor init multipart: bucket=%s, objectKey=%s", bucketName, objectKey)
		return "", nil
	}

	return uploadID, nil
}

func (u *usecaseImpl) CompleteUploadPart(bucketName string, objectKey string, uploadId string, parts []model.ObjectCompletePart) error {
	if err := u.storage.CompleteUploadPart(bucketName, objectKey, uploadId, parts); err != nil {
		log.WithError(err).Errorf("vendor complete multipart: bucket=%s, objectKey=%s, uploadId=%s, parts=%#v",
			bucketName, objectKey, uploadId, parts)
		return err
	}

	return nil
}

func (u *usecaseImpl) AbortMultipartUpload(bucketName string, objectKey string, uploadId string) error {
	if err := u.storage.AbortMultipartUpload(bucketName, objectKey, uploadId); err != nil {
		log.WithError(err).Errorf("vendor abort multipart: bucket=%s, objectKey=%s, uploadId=%s",
			bucketName, objectKey, uploadId)
		return err
	}

	return nil
}

func (u *usecaseImpl) UploadPart(bucketName string, objectKey string, uploadId string, partNum int, reader io.ReadSeeker) (string, error) {
	eTag, err := u.storage.UploadPart(bucketName, objectKey, uploadId, partNum, reader)
	if err != nil {
		log.WithError(err).Errorf("vendor upload multipart: bucket=%s, objectKey=%s, uploadId=%s, partNumber=%d",
			bucketName, objectKey, uploadId, partNum)
		return "", err
	}

	return eTag, nil
}

func (u *usecaseImpl) ListParts(bucketName string, objectKey string, uploadId string) ([]*model.ObjectPart, error) {
	parts, err := u.storage.ListParts(bucketName, objectKey, uploadId)
	if err != nil {
		log.WithError(err).Errorf("vendor list multipart: bucket=%s, objectKey=%s, uploadId=%s",
			bucketName, objectKey, uploadId)
		return nil, err
	}

	return parts, nil
}
