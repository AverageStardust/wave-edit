package wave

import (
	"io"
	"wave-edit/riff"
)

type DataChunk []byte

const dataChunkId = "data"

func dataChunkDeserializer(reader io.Reader, id riff.FourCC, size uint32) (DataChunk, error) {
	if id != dataChunkId {
		return nil, riff.ErrUnexpectedChunkId
	}

	data := make([]byte, size)
	n, err := reader.Read(data)
	if n < int(size) {
		return nil, riff.ErrUnexpectedEnd
	} else if err != nil && err != io.EOF {
		return nil, err
	}

	return DataChunk(data), nil
}

func (chunk DataChunk) Serialize(writer io.Writer) error {
	err := riff.SerializeChunkHeader(writer, dataChunkId, chunk.Size())
	if err != nil {
		return err
	}

	_, err = writer.Write(chunk)
	if err != nil {
		return err
	}

	return nil
}

func (chunk DataChunk) Size() uint32 {
	return 8 + uint32(len(chunk))
}
