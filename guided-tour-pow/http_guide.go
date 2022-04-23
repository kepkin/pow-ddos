package guidedTourPow

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

func ClientIdentity(r *http.Request) []byte {
	b := bytes.Buffer{}
	b.WriteString(r.RemoteAddr)
	b.WriteString(r.Header.Get("user-agent"))
	return b.Bytes()
}

type GuideRequest struct {
	Hash   Token
	Client string
	Step   int
	Length int
}

type HttpGuide struct {
	Endpoint string

	Client *http.Client
}

func NewHttpGuide(endpoint string) Guide {
	return &HttpGuide{
		Endpoint: endpoint,
		Client:   http.DefaultClient,
	}
}

func (g *HttpGuide) ComputeNextHash(ctx context.Context, previousHash Token, i int, client []byte, tourLength int) (Token, error) {

	body, err := json.Marshal(GuideRequest{
		Hash:   previousHash,
		Client: string(client),
		Step:   i,
		Length: tourLength,
	})

	if err != nil {
		return Token{}, err
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, g.Endpoint, bytes.NewBuffer(body))
	if err != nil {
		return Token{}, fmt.Errorf("failed to call remote guide: %w", err)
	}

	resp, err := g.Client.Do(httpReq)
	if err != nil {
		return Token{}, fmt.Errorf("failed to call remote guide: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return Token{}, fmt.Errorf("failed to call remote guide: status code %v", resp.Status)
	}

	var newHash Token
	dec := json.NewDecoder(resp.Body)
	err = dec.Decode(&newHash)
	if err != nil {
		return Token{}, err
	}

	return Token(newHash), nil
}
