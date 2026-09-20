package cacheSystem

import "time"

type Item struct {
	Value      []byte    `json:"value"`
	Expiration time.Time `json:"expiration"`
}

func (i Item) Expired() bool {

	if i.Expiration.IsZero() {
		return false
	}
	
	return time.Now().After(i.Expiration)
}
