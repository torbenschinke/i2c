package sht3x

import (
	"fmt"
	"github.com/torbenschinke/i2c/crc8"
	"io"
	"time"
)

var I2CAddrs = []uint16{
	0x44, // SHT3x default
	0x45, // SHT3x address B
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
	SOFTRESET    = 0x30A2
	READSTATUS   = 0xF32D
	CLEARSTATUS  = 0x3041
	MEAS_HIGHREP = 0x2400 // Measurement High Repeatability with Clock Stretch Disabled
)

func WriteCommand(w io.Writer, cmd Cmd) error {
	_, err := w.Write([]byte{byte(cmd >> 8), byte(cmd)})
	return err
}

func ClearStatus(w io.Writer) error {
	if err := WriteCommand(w, CLEARSTATUS); err != nil {
		return fmt.Errorf("cannot write CLEARSTATUS: %w", err)
	}

	time.Sleep(10 * time.Millisecond)
	return nil
}

func SoftReset(w io.Writer) error {
	if err := WriteCommand(w, SOFTRESET); err != nil {
		return fmt.Errorf("cannot write SOFTRESET: %w", err)
	}

	time.Sleep(300 * time.Millisecond)
	return nil
}

func ReadStatus(rw io.ReadWriter) (uint16, error) {
	v, err := readValue(rw, READSTATUS)
	return uint16(v), err
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

func readDoubleValue(rw io.ReadWriter, cmd Cmd) (int16, int16, error) {
	if err := WriteCommand(rw, cmd); err != nil {
		return 0, 0, fmt.Errorf("cannot write cmd %x: %w", cmd, err)
	}

	time.Sleep(20 * time.Millisecond)
	var buf [6]byte
	if _, err := rw.Read(buf[:]); err != nil {
		return 0, 0, fmt.Errorf("cannot read cmd %x: %w", cmd, err)
	}

	val0 := int16(buf[0])<<8 | int16(buf[1])
	if err := crc8.Checksum(crc8.Default, buf[:2], buf[2]); err != nil {
		return val0, 0, fmt.Errorf("invalid checksum on first tuple %v: %w", buf, err)
	}

	val1 := int16(buf[3])<<8 | int16(buf[4])
	if err := crc8.Checksum(crc8.Default, buf[3:5], buf[5]); err != nil {
		return val0, val1, fmt.Errorf("invalid checksum on second tuple %v: %w", buf, err)
	}

	return val0, val1, nil
}

func ReadTempHum(rw io.ReadWriter) (Temperature, Humidity, error) {
	temp, hum, err := readDoubleValue(rw, MEAS_HIGHREP)
	if err != nil {
		return 0, 0, fmt.Errorf("cannot read temp/hum data: %w", err)
	}

	return Temperature(temp), Humidity(hum), nil
}
