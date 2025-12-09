package riff

import (
	"errors"
	"io"
)

type Chunk interface {
	Serialize(writer io.Writer) error
	Size() uint32
}

type ChunkDeserializer[T Chunk] func(reader io.Reader, id FourCC, size uint32) (T, error)

var ErrUnexpectedEnd = errors.New("unexpected end of chunk")
var ErrReadTooMuch = errors.New("read more chunk content than expected")
var ErrUnexpectedChunkId = errors.New("unexpected chunk ID")

func SerializeChunkHeader(writer io.Writer, chunkId FourCC, size uint32) error {
	err := SerializeFourCC(writer, chunkId)
	if err != nil {
		return err
	}

	return SerializeDword(writer, size-8)
}

func DeserializeChunk[T Chunk](reader io.Reader, handler ChunkDeserializer[T]) (T, error) {
	var nothing T

	chunkId, err := DeserializeFourCC(reader)
	if err != nil {
		return nothing, err
	}

	dataSize, err := DeserializeDword(reader)
	if err != nil {
		return nothing, err
	}

	chunk, err := handler(reader, chunkId, dataSize)
	if err != nil {
		return nothing, err
	}

	return chunk, nil
}
