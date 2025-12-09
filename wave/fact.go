package wave

import (
	"io"
	"wave-edit/riff"
)

type FactChunk uint32

const factChunkId = "fact"

func factDeserializer(reader io.Reader, id riff.FourCC, _ uint32) (*FactChunk, error) {
	if id != factChunkId {
		return nil, riff.ErrUnexpectedChunkId
	}

	samples, err := riff.DeserializeDword(reader)
	if err != nil {
		return nil, err
	}

	chunk := FactChunk(samples)
	return &chunk, err
}

func (chunk *FactChunk) Serialize(writer io.Writer) error {
	err := riff.SerializeChunkHeader(writer, factChunkId, chunk.Size())
	if err != nil {
		return err
	}

	err = riff.SerializeDword(writer, uint32(*chunk))
	if err != nil {
		return err
	}

	return nil
}

func (chunk *FactChunk) Size() uint32 {
	return 12
}

func (chunk *FactChunk) Samples() uint32 {
	return uint32(*chunk)
}
