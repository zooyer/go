package fastdfs

import (
	"runtime"
	"fmt"
	"strings"
	"bytes"
	"errors"
)

const IGNORE = -1
const PAD    = -2
//const debug  = true

type Base64 struct {
	lineSeparator string
	lineLength    int
	valueToChar   [64]byte
	charToValue   [256]int
	charToPad     [4]byte
}

/* initialise defaultValueToChar and defaultCharToValue tables */
func initBase64(chPlus, chSplash, chPad byte) *Base64 {
	var this = new(Base64)
	var index = 0

	// build translate this.valueToChar table only once.
	// 0..25 -> 'A'..'Z'
	for i := 'A'; i <= 'Z'; i++ {
		this.valueToChar[index] = byte(i)
		index++
	}

	// 26..51 -> 'a'..'z'
	for i := 'a'; i <= 'z'; i++ {
		this.valueToChar[index] = byte(i)
		index++
	}

	// 52..61 -> '0'..'9'
	for i := '0'; i <= '9'; i++ {
		this.valueToChar[index] = byte(i)
		index++
	}

	this.valueToChar[index] = chPlus
	index++
	this.valueToChar[index] = chSplash
	index++

	// build translate defaultCharToValue table only once.
	for i := 0; i < 256; i++ {
		this.charToValue[i] = IGNORE  // default is to ignore
	}

	for i := 0; i < 64; i++ {
		this.charToValue[this.valueToChar[i]] = i
	}

	this.charToValue[chPad] = PAD

	for i,_ := range this.charToPad {
		this.charToPad[i] = chPad
	}

	return this
}

func NewBase64() *Base64 {
	return initBase64('+', '/', '=')
}

func NewBase64ByDetailed(chPlus, chSplash, chPad byte, lineLength int) *Base64 {
	var this = initBase64(chPlus, chSplash, chPad)
	this.lineLength = lineLength

	return this
}

func NewBase64ByLength(lineLength int) *Base64 {
	var this = new(Base64)
	this.lineLength = lineLength

	return this
}

/**
 * debug display array
 */
func show(b []byte) {
	var count = 0
	//var rows = 0

	for i := 0; i < len(b); i++ {
		if count == 8 {
			fmt.Print("  ")
		} else if count == 16 {
			fmt.Println("")
			count = 0
			continue
		}
		fmt.Print(strings.ToUpper(fmt.Sprintf("%02X", b[i])) + " ")
		count++
	}
	fmt.Println()
}

func display(b []byte) {
	fmt.Println(string(b))
}

func newLine() []byte {
	const (
		CR = '\r'
		LF = '\n'
	)
	switch runtime.GOOS {
	case "windows":
		return []byte{CR, LF}
	case "linux":
		fallthrough
	default:
		return []byte{LF}
	}

	return []byte{LF}
}

/**
 * Encode an arbitrary array of bytes as Base64 printable ASCII.
 * It will be broken into lines of 72 chars each.  The last line is not
 * terminated with a line separator.
 * The output will always have an even multiple of data characters,
 * exclusive of \n.  It is padded out with =.
 */
func (this *Base64) Encode(b []byte) (string, error) {
	// Each group or partial group of 3 bytes becomes four chars
	// covered quotient
	var outputLength = ((len(b) + 2) / 3) * 4

	// account for trailing newlines, on all but the very last line
	if this.lineLength != 0 {
		var lines = (outputLength + this.lineLength - 1) / this.lineLength - 1
		if lines > 0 {
			outputLength += lines * len(this.lineSeparator)
		}
	}

	// must be local for recursion to work.
	var sb = bytes.NewBuffer(nil)

	// must be local for recursion to work.
	var linePos = 0

	// first deal with even multiples of 3 bytes.
	var length = (len(b) / 3) * 3
	var leftover = len(b) - length

	for i := 0; i<length; i += 3 {
		// Start a new line if next 4 chars won't fit on the current line
		// We can't encapsulete the following code since the variable need to
		// be local to this incarnation of encode.
		linePos += 4
		if linePos > this.lineLength {
			if this.lineLength != 0 {
				if _,err := sb.WriteString(this.lineSeparator); err != nil {
					return "", err
				}
			}
			linePos = 4
		}

		// get next three bytes in unsigned form lined up,
		// in big-endian order
		var combined = b[i + 0] & 0xff
		combined <<= 8
		combined |= b[i + 1] & 0xff
		combined <<= 8
		combined |= b[i + 2] & 0xff

		// break those 24 bits into a 4 groups of 6 bits,
		// working LSB to MSB.
		var c3 = combined & 0x3f
		combined >>= 6
		var c2 = combined & 0x3f
		combined >>= 6
		var c1 = combined & 0x3f
		combined >>= 6
		var c0 = combined & 0x3f

		// Translate into the equivalent alpha character
		// emitting them in big-endian order.
		sb.WriteByte(this.valueToChar[c0])
		sb.WriteByte(this.valueToChar[c1])
		sb.WriteByte(this.valueToChar[c2])
		sb.WriteByte(this.valueToChar[c3])

		// deal with leftover bytes
		switch leftover {
		case 0:
			fallthrough
		default:
			// nothing to do
		case 1: {
			// One leftover byte generates xx==
			// Start a new line if next 4 chars won't fit on the current line
			linePos += 4
			if linePos > this.lineLength {

				if this.lineLength != 0 {
					if _,err := sb.WriteString(this.lineSeparator); err != nil {
						return "", err
					}
				}
				linePos = 4
			}

			// Handle this recursively with a faked complete triple.
			// Throw away last two chars and replace with ==
			str,err := this.Encode([]byte{b[length], 0, 0})
			if err != nil {
				return "", err
			}
			if _,err = sb.Write([]byte(str)[:2]); err != nil {
				return "", err
			}
			if _,err = sb.WriteString("=="); err != nil {
				return "", err
			}
		}
		case 2: {
			// Two leftover bytes generates xxx=
			// Start a new line if next 4 chars won't fit on the current line
			linePos += 4
			if linePos > this.lineLength {
				if this.lineLength != 0 {
					if _,err := sb.WriteString(this.lineSeparator); err != nil {
						return "", err
					}
				}
				linePos = 4
			}
			// Handle this recursively with a faked complete triple.
			// Throw away last char and replace with =
			str,err := this.Encode([]byte{b[length], b[length + 1], 0})
			if err != nil {
				return "", err
			}
			if _,err = sb.Write([]byte(str)[:3]); err != nil {
				return "", err
			}
			if _,err = sb.WriteString("="); err != nil {
				return "", err
			}
		}
		} // end switch
	}

	if outputLength != sb.Len() {
		fmt.Println("oops: minor program flaw: output length mis-estimated")
		fmt.Println("estimate:", outputLength)
		fmt.Println("actual:", sb.Len())
	}

	return sb.String(), nil
}

/**
 * decode a well-formed complete Base64 string back into an array of bytes.
 * It must have an even multiple of 4 data characters (not counting \n),
 * padded out with = as needed.
 */
func (this *Base64) DecodeAuto(s string) ([]byte, error) {
	var nRemain = len(s) % 4
	if nRemain == 0 {
		return this.Decode(s)
	}

	return this.Decode(s + string(this.charToPad[:4 - nRemain]))
}

/**
 * decode a well-formed complete Base64 string back into an array of bytes.
 * It must have an even multiple of 4 data characters (not counting \n),
 * padded out with = as needed.
 */
func (this *Base64) Decode(s string) ([]byte, error) {
	// estimate worst case size of output array, no embedded newlines.
	var b = make([]byte, (len(s) / 4) * 3)

	// tracks where we are in a cycle of 4 input chars.
	var cycle = 0

	// where we combine 4 groups of 6 bits and take apart as 3 groups of 8.
	var combined = 0

	// how many bytes we have prepared.
	var j = 0
	// will be an even multiple of 4 chars, plus some embedded \n
	var length = len(s)
	var dummies = 0
	for i := 0; i < length; i++ {
		var c = s[i]
		var value int
		if c <= 255 {
			value = this.charToValue[c]
		} else {
			value = IGNORE
		}
		// there are two magic values PAD (=) and IGNORE.
		switch value {
		case IGNORE:
			// e.g. \n, just ignore it.
		case PAD:
			value = 0
			dummies++
			fallthrough
		default:
			/* regular value character */
			switch cycle {
			case 0:
				combined = value
				cycle = 1
			case 1:
				combined <<= 6
				combined |= value
				cycle = 2
			case 2:
				combined <<= 6
				combined |= value
				cycle = 3
			case 3:
				combined <<= 6
				combined |= value
				// we have just completed a cycle of 4 chars.
				// the four 6-bit values are in combined in big-endian order
				// peel them off 8 bits at a time working lsb to msb
				// to get our original 3 8-bit bytes back

				b[j + 2] = byte(combined)
				combined >>= 8
				b[j + 1] = byte(combined)
				combined >>= 8
				b[j] = byte(combined)
				j += 3
				cycle = 0
			}
		}
	}

	if cycle != 0 {
		return nil, errors.New("input to decode not an even multiple of 4 characters; pad with =")
	}

	j -= dummies
	if len(b) != j {
		var b2 = make([]byte, j)
		copy(b2, b[:j])
		b = b2
	}

	return b, nil
}

/**
 * determines how long the lines are that are generated by encode.
 * Ignored by decode.
 *
 * @param length 0 means no newlines inserted. Must be a multiple of 4.
 */
func (this *Base64) SetLineLength(lineLength int) {
	this.lineLength = (lineLength / 4) * 4
}

/**
 * How lines are separated.
 * Ignored by decode.
 *
 * @param lineSeparator may be "" but not null.
 *                      Usually contains only a combination of chars \n and \r.
 *                      Could be any chars not in set A-Z a-z 0-9 + /.
 */
 func (this *Base64) SetLineSeparator(lineSeparator string) {
 	this.lineSeparator = lineSeparator
 }