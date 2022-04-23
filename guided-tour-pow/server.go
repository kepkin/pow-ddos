package guidedTourPow

import (
	"bytes"
	"context"
	"crypto/sha1"
	"encoding/binary"
	"fmt"
	"time"
)

func coarseTs() []byte {
	r, _ := time.Now().Truncate(time.Minute).MarshalBinary()
	return r
}

type Config struct {
	N       int      // Number of Guides
	Secrets [][]byte //TODO: token provider

	GuideID int // 0 stands for master server, which means you must have all secrets

	GetTS func() []byte
}

type Server interface {
	InitialHash(client []byte, tourLength int) Token
	Validate(h0 Token, hl Token, client []byte, tourLength int) (bool, error)
}

type Guide interface {
	ComputeNextHash(ctx context.Context, previousHash Token, i int, client []byte, tourLength int) (Token, error)
}

type impl struct {
	c Config
}

func NewServer(c Config) (Server, error) {
	if c.GetTS == nil {
		c.GetTS = coarseTs
	}

	if len(c.Secrets) != c.N+1 {
		return nil, fmt.Errorf("not all secrets provided")
	}

	if c.N < 2 {
		return nil, fmt.Errorf("minimal number of guides is 2 (%v provided)", c.N)
	}

	return &impl{
		c: c,
	}, nil
}

func NewGuide(guideID int, secret []byte, getTS func() []byte) (Guide, error) {
	if getTS == nil {
		getTS = coarseTs
	}

	if guideID == 0 {
		return nil, fmt.Errorf("guideID 0 is reserved for server")
	}

	return &impl{
		c: Config{
			Secrets: [][]byte{secret},
			GetTS:   getTS,
			GuideID: guideID,
		},
	}, nil
}

func (t impl) InitialHash(client []byte, tourLength int) Token {
	h := sha1.New()
	h.Write(client)
	binary.Write(h, binary.LittleEndian, tourLength)
	h.Write(t.c.GetTS())
	h.Write(t.c.Secrets[0])

	return Token(h.Sum(nil))
}

func (t impl) ComputeNextHash(ctx context.Context, previousHash Token, i int, client []byte, tourLength int) (Token, error) {
	if i > tourLength || i < 1 {
		return Token{}, fmt.Errorf("tour stop %v out of tour length %v", i, tourLength)
	}

	guideIndex := 0
	if t.c.GuideID == 0 { // Means it's a server with all secrets
		guideIndex = previousHash.Mod(t.c.N) + 1
	}

	if len(t.c.Secrets) < guideIndex || t.c.Secrets[guideIndex] == nil {
		return Token{}, fmt.Errorf("there is no secret for tour stop %v", i)
	}

	h := sha1.New()
	h.Write(previousHash.Bytes())
	binary.Write(h, binary.LittleEndian, i)
	binary.Write(h, binary.LittleEndian, tourLength)
	h.Write(client)
	h.Write(t.c.GetTS())
	h.Write(t.c.Secrets[guideIndex])

	return Token(h.Sum(nil)), nil
}

func (t impl) ComputeWholeTour(h0 Token, client []byte, tourLength int) (Token, error) {
	var err error
	for i := 1; i <= tourLength; i++ {
		h0, err = t.ComputeNextHash(context.Background(), h0, i, client, tourLength)
		if err != nil {
			return nil, err
		}
	}

	return h0, nil
}

func (t impl) Validate(h0 Token, hl Token, client []byte, tourLength int) (bool, error) {
	validHl, err := t.ComputeWholeTour(h0, client, tourLength)
	if err != nil {
		return false, err
	}

	return bytes.Equal(validHl, hl), nil
}
