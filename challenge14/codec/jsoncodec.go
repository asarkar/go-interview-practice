package codec

import (
	"encoding/json"

	"google.golang.org/grpc/encoding"
)

type jsonCodec struct{}

func (jsonCodec) Name() string {
	return "json"
}

func (jsonCodec) Marshal(v any) ([]byte, error) {
	return json.Marshal(v)
}

func (jsonCodec) Unmarshal(data []byte, v any) error {
	return json.Unmarshal(data, v)
}

func JSONCodec() encoding.Codec {
	var codec encoding.Codec
	if c := encoding.GetCodec("json"); c != nil {
		codec = c
	} else {
		codec = jsonCodec{}
		encoding.RegisterCodec(codec)
	}
	return codec
}
