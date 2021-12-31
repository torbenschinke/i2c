package shtc3

import (
	"fmt"
	"github.com/torbenschinke/i2c/crc8"
	"io"
	"time"
)

const I2CAddr = 0x70

type Ident uint16

// SHTC3 returns true, if Ident refers to the according sensor identifier.
// xxxx' 1 xxx’xx 00’0111
func (i Ident) SHTC3() bool {
	const (
		mask  = 0b0000100000111111
		ident = 0b0000100000000111
	)
	return i&mask == ident
}

// Cmd defines a 2 byte SHTC3 control opcode.
type Cmd uint16

type Temperature uint16

func (t Temperature) Celsius() float64 {
	return 175*float64(t)/65536.0 - 45.0
}

type Humidity uint16

func (h Humidity) Relative() float64 {
	return 100 * float64(h) / 65536.0
}

const (
	SHTC3_WakeUp Cmd = 0x3517
	Sleep        Cmd = 0xB098
	NM_CE_ReadTH Cmd = 0x7CA2
	NM_CE_ReadRH Cmd = 0x5C24
	NM_CD_ReadTH Cmd = 0x7866
	NM_CD_ReadRH Cmd = 0x58E0
	LM_CE_ReadTH Cmd = 0x6458
	LM_CE_ReadRH Cmd = 0x44DE
	LM_CD_ReadTH Cmd = 0x609C
	LM_CD_ReadRH Cmd = 0x401A
	Software_RES Cmd = 0x401A
	ID           Cmd = 0xEFC8
)

func WriteCommand(w io.Writer, cmd Cmd) error {
	_, err := w.Write([]byte{byte(cmd >> 8), byte(cmd)})
	return err
}

func SoftReset(w io.Writer) error {
	if err := WriteCommand(w, Software_RES); err != nil {
		return fmt.Errorf("cannot write Software_RES: %w", err)
	}

	time.Sleep(300 * time.Millisecond)
	return nil
}

func readValue(rw io.ReadWriter, cmd Cmd) (int16, error) {
	if err := WriteCommand(rw, cmd); err != nil {
		return 0, fmt.Errorf("cannot write cmd %x: %w", cmd, err)
	}

	time.Sleep(20 * time.Millisecond)
	var buf [3]byte
	if _, err := rw.Read(buf[:]); err != nil {
		return 0, fmt.Errorf("cannot read cmd %x: %w", cmd, err)
	}

	val := int16(buf[0])<<8 | int16(buf[1])
	if err := crc8.Checksum(crc8.Default, buf[:2], buf[2]); err != nil {
		fmt.Printf("invalid checksum %v\n", buf)
		return val, err
	}

	return val, nil
}

func ReadTemperature(rw io.ReadWriter) (Temperature, error) {
	val, err := readValue(rw, NM_CD_ReadTH)
	if err != nil {
		return 0, fmt.Errorf("cannot read temp: %w", err)
	}

	return Temperature(val), nil
}

func ReadHumidity(rw io.ReadWriter) (Humidity, error) {
	val, err := readValue(rw, NM_CD_ReadRH)
	if err != nil {
		return 0, fmt.Errorf("cannot read hum: %w", err)
	}

	return Humidity(val), nil
}

func ReadID(rw io.ReadWriter) (Ident, error) {
	val, err := readValue(rw, ID)
	if err != nil {
		return 0, fmt.Errorf("cannot read id: %w", err)
	}

	return Ident(val), nil
}
