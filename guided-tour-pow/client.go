package guidedTourPow

import (
	"context"
)

type Client struct {
	Guides []Guide
}

func (c Client) Pow(ctx context.Context, client []byte, h0 Token, tourLength int) (Token, error) {
	var err error
	ready := make(chan Token)

	go func() {
		hl := h0
		defer func() {

		}()

		for step := 1; step <= tourLength; step++ {
			g := c.Guides[hl.Mod(len(c.Guides))]
			hl, err = g.ComputeNextHash(ctx, hl, step, client, tourLength)
			if err != nil {
				ready <- hl
				return
			}
		}

		ready <- hl

	}()

	select {
	case <-ctx.Done():
		return Token{}, ctx.Err()
	case result := <-ready:
		return result, err
	}
}
