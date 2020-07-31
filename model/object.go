package model

import (
	obsSdk "github.com/cloud-mesh/object-storage-sdk"
	"time"
)

type Object struct {
	Vendor        string
	Bucket        string
	ObjectKey     string
	ContentType   string
	ContentLength int
	ETag          string
	LastModified  time.Time
}

type ObjectPart struct {
	Vendor       string
	Bucket       string
	ObjectKey    string
	PartNumber   int
	Size         int
	ETag         string
	LastModified time.Time
}

func AdapterObjectPart(vendor string, bucket string, objectKey string, part *obsSdk.Part) *ObjectPart {
	if part == nil {
		return nil
	}
	return &ObjectPart{
		Vendor:       vendor,
		Bucket:       bucket,
		ObjectKey:    objectKey,
		PartNumber:   part.PartNumber,
		Size:         part.Size,
		ETag:         part.ETag,
		LastModified: part.LastModified,
	}
}

func AdapterObjectParts(vendor string, bucket string, objectKey string, parts []obsSdk.Part) []*ObjectPart {
	objectParts := make([]*ObjectPart, 0, len(parts))
	for _, part := range parts {
		objectParts = append(objectParts, AdapterObjectPart(vendor, bucket, objectKey, &part))
	}

	return objectParts
}

func AdapterParts(objectParts []*ObjectPart) []obsSdk.Part {
	parts := make([]obsSdk.Part, 0, len(objectParts))
	for _, objectPart := range objectParts {
		parts = append(parts, obsSdk.Part{
			PartNumber:   objectPart.PartNumber,
			Size:         objectPart.Size,
			ETag:         objectPart.ETag,
			LastModified: objectPart.LastModified,
		})
	}

	return parts
}

type ObjectCompletePart struct {
	Vendor     string
	Bucket     string
	ObjectKey  string
	PartNumber int
	ETag       string
}

func AdapterObjectCompletePart(vendor string, bucket string, objectKey string, part *obsSdk.CompletePart) *ObjectCompletePart {
	if part == nil {
		return nil
	}
	return &ObjectCompletePart{
		Vendor:     vendor,
		Bucket:     bucket,
		ObjectKey:  objectKey,
		PartNumber: part.PartNumber,
		ETag:       part.ETag,
	}
}

func AdapterObjectCompleteParts(vendor string, bucket string, objectKey string, parts []obsSdk.CompletePart) []*ObjectCompletePart {
	completeParts := make([]*ObjectCompletePart, 0, len(parts))
	for _, part := range parts {
		completeParts = append(completeParts, AdapterObjectCompletePart(&part))
	}

	return completeParts
}

func AdapterCompleteParts(objectCompleteParts []*ObjectCompletePart) []obsSdk.CompletePart {
	parts := make([]obsSdk.CompletePart, 0, len(objectCompleteParts))
	for _, objectPart := range objectCompleteParts {
		parts = append(parts, obsSdk.CompletePart{
			PartNumber: objectPart.PartNumber,
			ETag:       objectPart.ETag,
		})
	}

	return parts
}
