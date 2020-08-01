package model

import (
	"time"
)

type Object struct {
	ObjectKey     string
	ContentType   string
	ContentLength int
	ETag          string
	LastModified  time.Time
}

type ObjectPart struct {
	PartNumber   int
	Size         int
	ETag         string
	LastModified time.Time
}

type ObjectCompletePart struct {
	PartNumber int
	ETag       string
}