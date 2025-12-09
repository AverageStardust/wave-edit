package riff

import "io"

type IgnoreChunk uint32

func IgnoreDeserializer(reader io.Reader, _ FourCC, size uint32) (Chunk, error) {
	n, err := reader.Read(make([]byte, size))

	if n < int(size) {
		return nil, ErrUnexpectedEnd
	} else if err != nil && err != io.EOF {
		return nil, err
	}

	return IgnoreChunk(size), nil
}

func (chunk IgnoreChunk) Serialize(writer io.Writer) error {
	return nil
}

func (chunk IgnoreChunk) Size() uint32 {
	return uint32(chunk) + 8
}
