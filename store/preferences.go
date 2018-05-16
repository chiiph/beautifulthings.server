package store

import "encoding/json"

const (
	Never  = "never"
	Daily  = "daily"
	Weekly = "weekly"
)

type Prefernces struct {
	nagInterval string `json:"nag-interval"`
}

func (p *Prefernces) Bytes() ([]byte, error) {
	return json.Marshal(p)
}
