package crc8

const InvalidChecksum constErr = "invalid crc checksum"

type constErr string

func (e constErr) Error() string {
	return string(e)
}

type Polynomial int16

// Default is the default Sensirion polynomial P(x) = x^8 + x^5 + x^4 + 1 = 100110001.
const Default = 0x131

func Checksum(poly Polynomial, data []byte, checksum byte) error {
	var bit byte        // bit mask
	var crc byte = 0xFF // calculated checksum
	// calculates 8-Bit checksum with given polynomial
	for _, b := range data {
		crc ^= b
		for bit = 8; int(bit) > 0; bit-- {
			if int(crc)&128 != 0 {
				crc = uint8((int(crc) << 1) ^ int(poly))
			} else {
				crc = uint8(int(crc) << 1)
			}
		}
	}

	if crc != checksum {
		return InvalidChecksum
	}

	return nil
}
