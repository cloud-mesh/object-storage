package http

import (
	"github.com/cloud-mesh/object-storage/model"
	"github.com/cloud-mesh/object-storage/utils/types"
)

type Bucket struct {
	Name      string         `json:"name"`
	Vendor    string         `json:"vendor"`
	CreatedAt types.JSONTime `json:"created_at"`
}

func adapterBucket(bucket *model.Bucket) *Bucket {
	if bucket == nil {
		return nil
	}
	return &Bucket{
		Name:      bucket.Name,
		Vendor:    bucket.Vendor,
		CreatedAt: types.JSONTime(bucket.CreatedAt),
	}
}

func adapterBuckets(buckets []*model.Bucket) []*Bucket {
	bucketViews := make([]*Bucket, 0, len(buckets))
	for _, application := range buckets {
		bucketView := adapterBucket(application)
		bucketViews = append(bucketViews, bucketView)
	}

	return bucketViews
}

type Object struct {
	Vendor        string         `json:"vendor"`
	Bucket        string         `json:"bucket"`
	ObjectKey     string         `json:"object_key"`
	ContentType   string         `json:"content_type"`
	ContentLength int            `json:"content_length"`
	ETag          string         `json:"etag"`
	LastModified  types.JSONTime `json:"last_modified"`
}

func adapterObject(object *model.Object) *Object {
	if object == nil {
		return nil
	}
	return &Object{
		Vendor:        object.Vendor,
		Bucket:        object.Bucket,
		ObjectKey:     object.ObjectKey,
		ContentType:   object.ContentType,
		ContentLength: object.ContentLength,
		ETag:          object.ETag,
		LastModified:  types.JSONTime(object.LastModified),
	}
}

func adapterObjects(objects []*model.Object) []*Object {
	objectViews := make([]*Object, 0, len(objects))
	for _, object := range objects {
		objectView := adapterObject(object)
		objectViews = append(objectViews, objectView)
	}

	return objectViews
}

type ObjectPart struct {
	Vendor       string         `json:"vendor"`
	Bucket       string         `json:"bucket"`
	ObjectKey    string         `json:"object_key"`
	PartNumber   int            `json:"part_number"`
	ETag         string         `json:"etag"`
	Size         int            `json:"size,omitempty"`
	LastModified types.JSONTime `json:"last_modified,omitempty"`
}

func adapterObjectPart(objectPart *model.ObjectPart) *ObjectPart {
	if objectPart == nil {
		return nil
	}
	return &ObjectPart{
		Vendor:       objectPart.Vendor,
		Bucket:       objectPart.Bucket,
		ObjectKey:    objectPart.ObjectKey,
		PartNumber:   objectPart.PartNumber,
		ETag:         objectPart.ETag,
		Size:         objectPart.Size,
		LastModified: types.JSONTime(objectPart.LastModified),
	}
}

func adapterObjectParts(objects []*model.ObjectPart) []*ObjectPart {
	partViews := make([]*ObjectPart, 0, len(objects))
	for _, object := range objects {
		partView := adapterObjectPart(object)
		partViews = append(partViews, partView)
	}

	return partViews
}
