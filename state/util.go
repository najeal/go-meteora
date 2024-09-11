package state

import (
	"encoding/base64"

	"github.com/near/borsh-go"
)

// DeserialieData expects a base64 encoded data
// it decodes and deserializes
func DeserializeData[T any](data string) (deserializedData T, err error) {
	decoded, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		return deserializedData, err
	}
	err = borsh.Deserialize(&deserializedData, decoded)
	if err != nil {
		return deserializedData, err
	}
	return deserializedData, nil
}

// SerialieData expects a base64 encoded data
// it decodes and deserializes
func SerializeData[T any](data T) (encodedData string, err error) {
	serializedData, err := borsh.Serialize(data)
	if err != nil {
		return encodedData, err
	}
	encodedData = base64.StdEncoding.EncodeToString(serializedData)
	if err != nil {
		return encodedData, err
	}
	return encodedData, nil
}
