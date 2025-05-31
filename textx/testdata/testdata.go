package testdata

import (
	"encoding/binary"
	"errors"
	"strconv"
	"time"
)

type Duration time.Duration

func (d *Duration) MarshalText() ([]byte, error) {
	return []byte(strconv.FormatInt(int64(*d), 10)), nil
}

func (d *Duration) UnmarshalText(data []byte) error {
	v, err := strconv.ParseInt(string(data), 10, 64)
	if err != nil {
		return err
	}
	*d = Duration(v)
	return nil
}

type Integers []uint64

func (d Integers) MarshalText() ([]byte, error) {
	data := make([]byte, len(d)*8)

	for i, v := range d {
		binary.LittleEndian.PutUint64(data[i*8:], v)
	}
	return data, nil
}

func (d *Integers) UnmarshalText(data []byte) error {
	*d = make(Integers, len(data)/8)
	if len(data) == 0 {
		return nil
	}
	if len(data)%8 != 0 {
		return errors.New("invalid length of input")
	}
	for offset := 0; len(data) >= 8; offset += 8 {
		v := binary.LittleEndian.Uint64(data)
		data = data[8:]
		(*d)[offset/8] = v
	}
	return nil
}

type MustFailedArshaler struct {
	V any
}

func (MustFailedArshaler) MarshalText() ([]byte, error) {
	return nil, errors.New("")
}

func (*MustFailedArshaler) UnmarshalText([]byte) error {
	return errors.New("")
}

type (
	String string
	Bytes  []byte
)
