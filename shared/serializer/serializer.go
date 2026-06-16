package serializer

import "github.com/bytedance/sonic"

func Marshal(value any) ([]byte, error) {
	return sonic.Marshal(value)
}

func Unmarshal(data []byte, value any) error {
	return sonic.Unmarshal(data, value)
}
