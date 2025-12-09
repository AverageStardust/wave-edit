package riff

import (
	"errors"
	"io"
)

type FormDeserializer func(reader io.Reader, size uint32) (Chunk, error)

var ErrNotRIFF = errors.New("not a RIFF file")
var ErrUnexpectedForm = errors.New("unexpected RIFF form")

var formDeserializers = map[FourCC]FormDeserializer{}

func SerializeRiffHeader(writer io.Writer, size uint32, form FourCC) error {
	err := SerializeChunkHeader(writer, "RIFF", size)
	if err != nil {
		return err
	}

	err = SerializeFourCC(writer, form)
	if err != nil {
		return err
	}

	return nil
}

func DeserializerRiff(reader io.Reader) (Chunk, error) {
	return DeserializeChunk(reader,
		func(reader io.Reader, id FourCC, size uint32) (Chunk, error) {
			if id != "RIFF" {
				return nil, ErrNotRIFF
			}

			form, err := DeserializeFourCC(reader)
			size -= 4
			if err != nil {
				return nil, err
			}

			deserializer := formDeserializers[form]
			if deserializer != nil {
				return deserializer(reader, size)
			} else {
				return nil, ErrUnexpectedForm
			}
		})
}

func RegisterRiffForm(form FourCC, deserializer FormDeserializer) {
	formDeserializers[form] = deserializer
}
