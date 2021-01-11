package hbhsid

import (
	"encoding/json"
)

func (id *ID) UnmarshalJSON(src []byte) error {
	var s string
	if err := json.Unmarshal(src, &s); err != nil {
		return err
	}
	return id.FromString(s)
}

func (id ID) MarshalJSON() ([]byte, error) {
	return json.Marshal(id.String())
}
