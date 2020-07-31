package model

import "time"

type Bucket struct {
	Vendor    string
	Name      string
	CreatedAt time.Time
}
