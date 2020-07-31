package storage

import (
	"errors"
	obsSdk "github.com/cloud-mesh/object-storage-sdk"
	"github.com/cloud-mesh/object-storage/model"
	"io"
)

type repoImpl struct {
	vendor string
	client obsSdk.BasicClient
}

func New(vendor string, client obsSdk.BasicClient) *repoImpl {
	return &repoImpl{vendor, client}
}

func (r *repoImpl) GetBucket(bucketName string) (*model.Bucket, error) {
	if err := r.client.HeadBucket(bucketName); err != nil {
		return nil, err
	}

	bucket := &model.Bucket{
		Vendor: r.vendor,
		Name:   bucketName,
	}
	return bucket, nil
}

func (r *repoImpl) ListBucket(page int, pageSize int) ([]*model.Bucket, error) {
	bucketProperties, err := r.client.ListBucket() // TODO: support page & pageSize
	if err != nil {
		return nil, err
	}

	return adapterBuckets(r.vendor, bucketProperties), nil
}

func (r *repoImpl) CreateBucket(bucketName string) error {
	return r.client.MakeBucket(bucketName)
}

func (r *repoImpl) DeleteBucket(bucketName string) error {
	return r.client.RemoveBucket(bucketName)
}

func (r *repoImpl) HeadObject(bucketName string, objectKey string) (object *model.Object, err error) {
	bucket, err := r.client.Bucket(bucketName)
	if err != nil {
		return nil, err
	}

	objectMeta, err := bucket.HeadObject(objectKey)
	if err != nil {
		return nil, err
	}

	return adapterObject(r.vendor, bucketName, objectKey, objectMeta), nil
}

func (r *repoImpl) PutObject(bucketName string, objectKey string, reader io.Reader) error {
	bucket, err := r.client.Bucket(bucketName)
	if err != nil {
		return err
	}

	return bucket.PutObject(objectKey, reader)
}

func (r *repoImpl) DeleteObject(bucketName string, objectKey string) error {
	bucket, err := r.client.Bucket(bucketName)
	if err != nil {
		return err
	}

	return bucket.RemoveObject(objectKey)
}

func (r *repoImpl) InitMultipartUpload(bucketName string, objectKey string) (uploadId string, err error) {
	bucket, err := r.getMultipartBucket(bucketName)
	if err != nil {
		return "", err
	}

	return bucket.InitMultipartUpload(objectKey)
}

func (r *repoImpl) CompleteUploadPart(bucketName string, objectKey string, uploadId string, parts []model.ObjectCompletePart) error {
	bucket, err := r.getMultipartBucket(bucketName)
	if err != nil {
		return err
	}

	return bucket.CompleteUploadPart(objectKey, uploadId, adapterCompleteParts(parts))
}

func (r *repoImpl) AbortMultipartUpload(bucketName string, objectKey string, uploadId string) error {
	bucket, err := r.getMultipartBucket(bucketName)
	if err != nil {
		return err
	}

	return bucket.AbortMultipartUpload(objectKey, uploadId)
}

func (r *repoImpl) UploadPart(bucketName string, objectKey string, uploadId string, partNum int, reader io.ReadSeeker) (string, error) {
	bucket, err := r.getMultipartBucket(bucketName)
	if err != nil {
		return "", err
	}

	return bucket.UploadPart(objectKey, uploadId, partNum, reader)
}

func (r *repoImpl) ListParts(bucketName string, objectKey string, uploadId string) ([]*model.ObjectPart, error) {
	bucket, err := r.getMultipartBucket(bucketName)
	if err != nil {
		return nil, err
	}

	parts, err := bucket.ListParts(objectKey, uploadId)
	if err != nil {
		return nil, err
	}

	return adapterObjectParts(r.vendor, bucketName, objectKey, parts), nil
}

func (r *repoImpl) getMultipartBucket(bucketName string) (obsSdk.MultipartUploadAbleBucket, error) {
	bucket, err := r.client.Bucket(bucketName)
	if err != nil {
		return nil, err
	}

	multipartBucket, ok := bucket.(obsSdk.MultipartUploadAbleBucket)
	if !ok {
		return nil, errors.New("vendor not support multipart upload")
	}

	return multipartBucket, nil
}
