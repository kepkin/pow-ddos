package guidedTourPow

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"
)

type Token []byte

func (t Token) Mod(m int) int {
	x, _ := binary.Uvarint([]byte(t))
	return int(x % uint64(m))
}

func (t Token) Bytes() []byte {
	return []byte(t)
}

func (t Token) Hex() string {
	return hex.EncodeToString([]byte(t))
}

func NewTokenFromHex(src string) Token {
	r, _ := hex.DecodeString(src)
	return r
}

func (t *Token) UnmarshalJSON(src []byte) error {
	if len(src) < 2 || src[0] != '"' {
		return fmt.Errorf("token json unmarshal requires string type")
	}
	newT := make([]byte, hex.DecodedLen(len(src)-2))
	_, err := hex.Decode(newT, src[1:len(src)-1])

	*t = Token(newT)
	return err
}

func (t Token) MarshalJSON() ([]byte, error) {
	return []byte("\"" + t.Hex() + "\""), nil
}
