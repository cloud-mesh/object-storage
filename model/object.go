package model

import (
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

type ObjectCompletePart struct {
	Vendor     string
	Bucket     string
	ObjectKey  string
	PartNumber int
	ETag       string
}