package storage

import (
	"errors"
	obsSdk "github.com/cloud-mesh/object-storage-sdk"
	"github.com/cloud-mesh/object-storage/model"
	"io"
)

type repoImpl struct {
	clientAdapter obsSdk.BasicClient
}

func New(clientAdapter obsSdk.BasicClient) *repoImpl {
	return &repoImpl{clientAdapter}
}

func (r *repoImpl) GetBucket(bucketName string) (*model.Bucket, error) {
	if err := r.clientAdapter.HeadBucket(bucketName); err != nil {
		return nil, err
	}

	bucket := &model.Bucket{
		Name: bucketName,
	}
	return bucket, nil
}

func (r *repoImpl) ListBucket(page int, pageSize int) ([]*model.Bucket, error) {
	bucketProperties, err := r.clientAdapter.ListBucket() // TODO: support page & pageSize
	if err != nil {
		return nil, err
	}

	return adapterBuckets(bucketProperties), nil
}

func (r *repoImpl) CreateBucket(bucketName string) error {
	return r.clientAdapter.MakeBucket(bucketName)
}

func (r *repoImpl) DeleteBucket(bucketName string) error {
	return r.clientAdapter.RemoveBucket(bucketName)
}

func (r *repoImpl) HeadObject(bucketName string, objectKey string) (object *model.Object, err error) {
	bucket, err := r.clientAdapter.Bucket(bucketName)
	if err != nil {
		return nil, err
	}

	objectMeta, err := bucket.HeadObject(objectKey)
	if err != nil {
		return nil, err
	}

	return adapterObject(objectKey, objectMeta), nil
}

func (r *repoImpl) PutObject(bucketName string, objectKey string, reader io.Reader) error {
	bucket, err := r.clientAdapter.Bucket(bucketName)
	if err != nil {
		return err
	}

	return bucket.PutObject(objectKey, reader)
}

func (r *repoImpl) DeleteObject(bucketName string, objectKey string) error {
	bucket, err := r.clientAdapter.Bucket(bucketName)
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

	return adapterObjectParts(parts), nil
}

func (r *repoImpl) getMultipartBucket(bucketName string) (obsSdk.MultipartUploadAbleBucket, error) {
	bucket, err := r.clientAdapter.Bucket(bucketName)
	if err != nil {
		return nil, err
	}

	multipartBucket, ok := bucket.(obsSdk.MultipartUploadAbleBucket)
	if !ok {
		return nil, errors.New("vendor not support multipart upload")
	}

	return multipartBucket, nil
}
