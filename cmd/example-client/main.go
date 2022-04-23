package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/alexflint/go-arg"

	pow "github.com/kepkin/pow-ddos/guided-tour-pow"
)

func makeServerReq(ctx context.Context, s string, powClient pow.Client) error {
	req, err := http.NewRequest("GET", s, nil)
	if err != nil {
		return err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized {
		h0 := resp.Header.Get("x-pow-hash0")
		client := resp.Header.Get("x-pow-client")
		l := resp.Header.Get("x-pow-length")

		tourLength, err := strconv.Atoi(l)
		hl, err := powClient.Pow(ctx, []byte(client), pow.NewTokenFromHex(h0), tourLength)
		if err != nil {
			return err
		}

		hreq, err := http.NewRequest("GET", s, nil)
		xPowHashes := strings.Join([]string{h0, hl.Hex(), l}, ",")
		hreq.Header.Add("x-pow-hashes", xPowHashes)

		hresp, err := http.DefaultClient.Do(hreq)
		if err != nil {
			return err
		}
		defer hresp.Body.Close()

		if hresp.StatusCode != http.StatusOK {
			return fmt.Errorf("failed")
		}

		data, err := ioutil.ReadAll(hresp.Body)
		if err != nil {
			return err
		}

		fmt.Println(string(data))
	}
	return nil
}

func main() {
	var args struct {
		Server string
		Guides []string
	}
	arg.MustParse(&args)

	gs := make([]pow.Guide, len(args.Guides))
	for i, g := range args.Guides {
		gs[i] = pow.NewHttpGuide(g)
	}

	powClient := pow.Client{gs}
	err := makeServerReq(context.TODO(), args.Server, powClient)
	if err != nil {
		log.Fatal(err)
	}
}
