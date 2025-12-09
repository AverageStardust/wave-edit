package riff

import (
	"encoding/binary"
	"errors"
	"io"
	"reflect"
)

type FourCC string

var ErrInvalidFourCC = errors.New("invalid four char code, wrong size")

func SerializeStruct(writer io.Writer, data any) error {
	value := reflect.ValueOf(data)

	for i := range value.NumField() {
		switch value.Field(i).Type() {
		case reflect.TypeFor[uint32]():
			err := SerializeDword(writer, uint32(value.Field(i).Uint()))
			if err != nil {
				return err
			}

		case reflect.TypeFor[uint16]():
			err := SerializeWord(writer, uint16(value.Field(i).Uint()))
			if err != nil {
				return err
			}

		case reflect.TypeFor[FourCC]():
			err := SerializeFourCC(writer, FourCC(value.Field(i).String()))
			if err != nil {
				return err
			}

		}
	}

	return nil
}

func DeserializeStruct[T any](reader io.Reader) (T, error) {
	value := reflect.New(reflect.TypeFor[T]()).Elem()

	var nothing T
	for i := range value.NumField() {
		switch value.Field(i).Type() {
		case reflect.TypeFor[uint32]():
			dword, err := DeserializeDword(reader)
			if err != nil {
				return nothing, err
			}

			value.Field(i).SetUint(uint64(dword))

		case reflect.TypeFor[uint16]():
			word, err := DeserializeWord(reader)
			if err != nil {
				return nothing, err
			}

			value.Field(i).SetUint(uint64(word))

		case reflect.TypeFor[FourCC]():
			fourCC, err := DeserializeFourCC(reader)
			if err != nil {
				return nothing, err
			}

			value.Field(i).SetString(string(fourCC))
		}
	}

	return value.Interface().(T), nil
}

func SerializeFourCC(writer io.Writer, code FourCC) error {
	if len(code) != 4 {
		return ErrInvalidFourCC
	}

	_, err := writer.Write([]byte(code))
	return err
}

func DeserializeFourCC(reader io.Reader) (FourCC, error) {
	var buffer [4]byte
	n, err := reader.Read(buffer[:])

	if n < 4 {
		return "", ErrUnexpectedEnd
	} else if err != nil && err != io.EOF {
		return "", err
	}

	return FourCC(buffer[:]), nil
}

func SerializeDword(writer io.Writer, word uint32) error {
	var buffer [4]byte
	binary.LittleEndian.PutUint32(buffer[:], word)

	_, err := writer.Write(buffer[:])

	return err
}

func DeserializeDword(reader io.Reader) (uint32, error) {
	var buffer [4]byte
	n, err := reader.Read(buffer[:])

	if n < 4 {
		return 0, ErrUnexpectedEnd
	} else if err != nil && err != io.EOF {
		return 0, err
	}

	return binary.LittleEndian.Uint32(buffer[:]), nil
}

func SerializeWord(writer io.Writer, word uint16) error {
	var buffer [2]byte
	binary.LittleEndian.PutUint16(buffer[:], word)

	_, err := writer.Write(buffer[:])

	return err
}

func DeserializeWord(reader io.Reader) (uint16, error) {
	var buffer [2]byte
	n, err := reader.Read(buffer[:])

	if n < 2 {
		return 0, ErrUnexpectedEnd
	} else if err != nil && err != io.EOF {
		return 0, err
	}

	return binary.LittleEndian.Uint16(buffer[:]), nil
}
