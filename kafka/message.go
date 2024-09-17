package kafka

import "encoding/json"

func NewJSONMessage(key string, data any) (Message, error) {
	jsonBytes, err := json.Marshal(data)
	if err != nil {
		return Message{}, err
	}

	return Message{
		Key:   []byte(key),
		Value: jsonBytes,
	}, nil
}
