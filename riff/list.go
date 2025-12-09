package riff

import (
	"errors"
	"io"
)

type ListChunk[T Chunk] struct {
	ListType FourCC
	Chunks   []T
}

const ListChunkId = "LIST"

var ErrUnexpectedListType = errors.New("unexpected list chunk type")

func ListChunkDeserializer[T Chunk](reader io.Reader, id FourCC, size uint32, expectedType FourCC, elementHandler ChunkDeserializer[T]) (*ListChunk[T], error) {
	if id != ListChunkId {
		return nil, ErrUnexpectedChunkId
	}

	listType, err := DeserializeFourCC(reader)
	size -= 4

	if err != nil {
		return nil, err
	} else if listType != expectedType {
		return nil, ErrUnexpectedListType
	}

	elements := []T{}
	for size > 0 {
		elm, err := DeserializeChunk(reader,
			func(reader io.Reader, elmId FourCC, elmSize uint32) (T, error) {
				if size < elmSize {
					var nothing T
					return nothing, ErrReadTooMuch
				}

				size -= elmSize
				return elementHandler(reader, elmId, elmSize)
			})

		if err != nil {
			return nil, err
		}
		elements = append(elements, elm)
	}

	return &ListChunk[T]{
		listType,
		elements,
	}, nil
}

func (chunk *ListChunk[T]) Serialize(writer io.Writer) error {
	err := SerializeChunkHeader(writer, ListChunkId, chunk.Size())
	if err != nil {
		return err
	}

	err = SerializeFourCC(writer, chunk.ListType)
	if err != nil {
		return err
	}

	for _, child := range chunk.Chunks {
		child.Serialize(writer)
	}

	return nil
}

func (chunk *ListChunk[T]) Size() uint32 {
	var size uint32 = 12

	for _, child := range chunk.Chunks {
		size += child.Size()
	}

	return size
}
