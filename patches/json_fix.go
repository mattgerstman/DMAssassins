
package uuid

import (
	"encoding/json"
)

func (uuid UUID) MarshalJSON() ([]byte, error) {
	return json.Marshal(uuid.String())
}

func (uuid *UUID) UnmarshalJSON(in []byte) error {
	var str string
	err := json.Unmarshal(in, &str)
	if err != nil {
		return err
	}
	*uuid = (*uuid)[:0]
	id := Parse(str)
	if id != nil {
		*uuid = append(*uuid, id...)
	} else {
		// return an error here
	}
	return nil
}
