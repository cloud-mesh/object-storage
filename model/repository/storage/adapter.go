package storage

import (
	obsSdk "github.com/cloud-mesh/object-storage-sdk"
	"github.com/cloud-mesh/object-storage/model"
)

func adapterBucket(bucketProperties obsSdk.BucketProperties) *model.Bucket {
	return &model.Bucket{
		Name:   bucketProperties.Name,
	}
}

func adapterBuckets(bucketProperties []obsSdk.BucketProperties) []*model.Bucket {
	buckets := make([]*model.Bucket, 0, len(bucketProperties))
	for _, bucketProperty := range bucketProperties {
		buckets = append(buckets, adapterBucket(bucketProperty))
	}

	return buckets
}

func adapterObject(objectKey string, objectMeta obsSdk.ObjectMeta) *model.Object {
	return &model.Object{
		ObjectKey:     objectKey,
		ContentType:   objectMeta.ContentType,
		ContentLength: objectMeta.ContentLength,
		ETag:          objectMeta.ETag,
		LastModified:  objectMeta.LastModified,
	}
}

func adapterObjectPart(part *obsSdk.Part) *model.ObjectPart {
	if part == nil {
		return nil
	}
	return &model.ObjectPart{
		PartNumber:   part.PartNumber,
		Size:         part.Size,
		ETag:         part.ETag,
		LastModified: part.LastModified,
	}
}

func adapterObjectParts(parts []obsSdk.Part) []*model.ObjectPart {
	objectParts := make([]*model.ObjectPart, 0, len(parts))
	for _, part := range parts {
		objectParts = append(objectParts, adapterObjectPart(&part))
	}

	return objectParts
}

func adapterCompleteParts(objectCompleteParts []model.ObjectCompletePart) []obsSdk.CompletePart {
	parts := make([]obsSdk.CompletePart, 0, len(objectCompleteParts))
	for _, objectPart := range objectCompleteParts {
		parts = append(parts, obsSdk.CompletePart{
			PartNumber: objectPart.PartNumber,
			ETag:       objectPart.ETag,
		})
	}

	return parts
}
