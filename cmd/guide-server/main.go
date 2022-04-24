package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/alexflint/go-arg"

	pow "github.com/kepkin/pow-ddos/guided-tour-pow"
)

type GuideHandler struct {
	g pow.Guide
}

func (h *GuideHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var guideReq pow.GuideRequest
	err := json.NewDecoder(r.Body).Decode(&guideReq)
	if err != nil {
		log.Print("failed to calc next hash: ", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	newHash, err := h.g.ComputeNextHash(r.Context(), guideReq.Hash, guideReq.Step, []byte(guideReq.Client), guideReq.Length)
	if err != nil {
		log.Print("failed to calc next hash: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	enc := json.NewEncoder(w)
	err = enc.Encode(newHash.Hex())
	if err != nil {
		log.Print("failed to make response: ", err)
	}
	log.Print("issued: ", newHash.Hex())
}

func main() {
	var args struct {
		Addr    string `default:":8080"`
		GuideID int
		Secret  string
	}
	arg.MustParse(&args)
	log.Print("Secret: ", args.Secret)

	g, err := pow.NewGuide(args.GuideID, []byte(args.Secret), nil)
	if err != nil {
		log.Fatal(err)
	}

	s := &http.Server{
		Addr:    args.Addr,
		Handler: &GuideHandler{g},
	}
	log.Fatal(s.ListenAndServe())
}
