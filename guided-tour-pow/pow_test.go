package guidedTourPow

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func mockGetTS() []byte {
	return []byte("timestamp")
}

func TestValidation(t *testing.T) {
	a := assert.New(t)

	c := Config{
		Secrets: [][]byte{
			[]byte("one"),
			[]byte("guide1"),
			[]byte("guide2"),
		},
		N:     2,
		GetTS: mockGetTS,
	}

	client := []byte("client1")

	server, err := NewServer(c)
	a.NoError(err)

	nextGuideIdx := 1
	nextGuide := func(k string) Guide {
		g, err := NewGuide(nextGuideIdx, []byte(k), mockGetTS)
		nextGuideIdx += 1
		a.NoError(err)

		return g
	}

	guides := []Guide{
		nextGuide("guide1"),
		nextGuide("guide2"),
	}

	tourLength := 10
	hashes := make([]Token, 0)
	hashes = append(hashes, server.InitialHash(client, tourLength))

	for step := 1; step <= tourLength; step++ {
		prevH := hashes[len(hashes)-1]
		g := guides[prevH.Mod(2)]
		hi, err := g.ComputeNextHash(context.Background(), prevH, step, client, tourLength)
		a.NoError(err)
		hashes = append(hashes, hi)
	}

	r, err := server.Validate(hashes[0], hashes[len(hashes)-1], client, tourLength)
	a.NoError(err)
	a.Truef(r, "won't validate h0: %x h1: %x", hashes[0], hashes[len(hashes)-1])

	if t.Failed() {
		t.Logf("list of hashes: %x", hashes)
	}
}
