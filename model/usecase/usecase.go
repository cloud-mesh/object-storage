package usecase

import (
	"github.com/cloud-mesh/object-storage/model"
	"io"
)

type UseCase interface {
	GetBucket(vendor string, bucketName string) (bucket *model.Bucket, err error)
	ListBucket(vendor string, page int, pageSize int) (buckets []*model.Bucket, err error)
	CreateBucket(vendor string, bucketName string) error
	DeleteBucket(vendor string, bucketName string) error

	HeadObject(vendor string, bucketName string, objectKey string) (object *model.Object, err error)
	PutObject(vendor string, bucketName string, objectKey string, reader io.Reader) error
	DeleteObject(vendor string, bucketName string, objectKey string) error

	InitMultipartUpload(vendor string, bucketName string, objectKey string) (uploadId string, err error)
	CompleteUploadPart(vendor string, bucketName string, objectKey string, uploadId string, parts []model.ObjectCompletePart) error
	AbortMultipartUpload(vendor string, bucketName string, objectKey string, uploadId string) error
	UploadPart(vendor string, bucketName string, objectKey string, uploadId string, partNum int, reader io.ReadSeeker) (eTag string, err error)
	ListParts(vendor string, bucketName string, objectKey string, uploadId string) (parts []*model.ObjectPart, err error)
}
