package main

import (
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"strings"

	"github.com/alexflint/go-arg"

	pow "github.com/kepkin/pow-ddos/guided-tour-pow"
)

type ExampleHandler struct {
	s pow.Server
}

func (h *ExampleHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	client := pow.ClientIdentity(r)

	powHashes := r.Header.Get("x-pow-hashes")
	if len(powHashes) == 0 {
		tourLength := 5
		h0 := h.s.InitialHash(client, tourLength)

		w.Header().Add("x-pow-hash0", h0.Hex())
		w.Header().Add("x-pow-length", strconv.Itoa(tourLength))
		w.Header().Add("x-pow-client", string(client))

		log.Print("issued: ", h0.Hex())

		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	hashes := strings.SplitN(powHashes, ",", 3)
	if len(hashes) != 3 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	tourLength, err := strconv.Atoi(hashes[2])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	valid, err := h.s.Validate(pow.NewTokenFromHex(hashes[0]), pow.NewTokenFromHex(hashes[1]), client, tourLength)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if !valid {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	w.WriteHeader(http.StatusOK)

	quote := quotes[rand.Intn(len(quotes))]
	w.Write([]byte(quote))
}

func main() {
	var args struct {
		Addr    string `default:":8080"`
		Secrets []string
	}
	arg.MustParse(&args)

	secrets := make([][]byte, len(args.Secrets))
	for i, v := range args.Secrets {
		secrets[i] = []byte(v)
	}

	log.Print(secrets)

	c := pow.Config{
		Secrets: secrets,
		N:       len(args.Secrets) - 1,
	}

	g, err := pow.NewServer(c)
	if err != nil {
		log.Fatal(err)
	}

	s := &http.Server{
		Addr:    args.Addr,
		Handler: &ExampleHandler{g},
	}
	log.Fatal(s.ListenAndServe())
}
