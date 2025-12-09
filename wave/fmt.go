package wave

import (
	"errors"
	"io"
	"wave-edit/riff"
)

type FmtChunk struct {
	Format        WaveFormat // Format
	Channels      uint16     // Number of channels
	SamplesPerSec uint32     // Sampling rate
}

type rawFmtChunk struct {
	FormatTag      uint16 // Format category
	Channels       uint16 // Number of channels
	SamplesPerSec  uint32 // Sampling rate
	AvgBytesPerSec uint32 // For buffer estimation
	BlockSize      uint16 // Data block size
	BitsPerSample  uint16 // Sample size
}

const fmtChunkId = "fmt "

var ErrUnsupportedFormat = errors.New("unsupported WAVE data format")

func fmtDeserializer(reader io.Reader, id riff.FourCC, _ uint32) (*FmtChunk, error) {
	if id != fmtChunkId {
		return nil, riff.ErrUnexpectedChunkId
	}

	rawFmt, err := riff.DeserializeStruct[rawFmtChunk](reader)
	if err != nil {
		return nil, err
	}

	format := createWaveFormat(rawFmt.FormatTag, rawFmt.BitsPerSample)

	if format == UNKNOWN_FORMAT {
		return nil, ErrUnsupportedFormat
	}

	return &FmtChunk{
		Format:        format,
		Channels:      rawFmt.Channels,
		SamplesPerSec: rawFmt.SamplesPerSec,
	}, nil
}

func (chunk *FmtChunk) Serialize(writer io.Writer) error {
	err := riff.SerializeChunkHeader(writer, fmtChunkId, chunk.Size())
	if err != nil {
		return err
	}

	formatTag, byteDepth := chunk.Format.Properties()
	blockSize := chunk.Channels * byteDepth

	rawFmt := rawFmtChunk{
		FormatTag:      formatTag,
		Channels:       chunk.Channels,
		SamplesPerSec:  chunk.SamplesPerSec,
		AvgBytesPerSec: uint32(blockSize) * chunk.SamplesPerSec,
		BlockSize:      blockSize,
		BitsPerSample:  byteDepth * 8,
	}

	err = riff.SerializeStruct(writer, rawFmt)
	if err != nil {
		return err
	}

	return nil
}

func (chunk *FmtChunk) Size() uint32 {
	return 8 + 16
}

func (chunk *FmtChunk) BlockSize() uint16 {
	_, byteDepth := chunk.Format.Properties()
	return chunk.Channels * byteDepth
}
