package storage

import (
	"github.com/cloud-mesh/object-storage/model"
	"io"
)

type Repo interface {
	GetBucket(bucketName string) (bucket *model.Bucket, err error)
	ListBucket(page int, pageSize int) (buckets []*model.Bucket, err error)
	CreateBucket(bucketName string) error
	DeleteBucket(bucketName string) error

	CreateObject(bucketName string, objectKey string, reader io.Reader) error
	HeadObject(bucketName string, objectKey string) (object *model.Object, err error)
	PutObject(bucketName string, objectKey string, reader io.Reader) error
	DeleteObject(bucketName string, objectKey string) error

	InitMultipartUpload(bucketName string, objectKey string) (uploadId string, err error)
	CompleteUploadPart(bucketName string, objectKey string, uploadId string, parts []*model.ObjectCompletePart) error
	AbortMultipartUpload(bucketName string, objectKey string, uploadId string) error
	UploadPart(bucketName string, objectKey string, uploadId string, partNum int, reader io.ReadSeeker) (eTag string, err error)
	ListParts(bucketName string, objectKey string, uploadId string) (parts []*model.ObjectPart, err error)
}
