package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	gitopssetsv1 "github.com/gitops-tools/gitopssets-controller/api/v1alpha1"
)

func newGeneratorMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("POST /generate", func(w http.ResponseWriter, req *http.Request) {
		var generatorRequest gitopssetsv1.GeneratorRequest
		decoder := json.NewDecoder(req.Body)
		if err := decoder.Decode(&generatorRequest); err != nil {
			log.Printf("failed to decode Generator Request: %s", err)
			http.Error(w, "unable to decode request", http.StatusBadRequest)
			return
		}

		wantUsers := int64(generatorRequest.GeneratorParams["count"].(float64))
		log.Printf("generating %v users", wantUsers)

		response := gitopssetsv1.GeneratorResponse{Generated: generateRandomUsers(wantUsers)}
		w.Header().Set("Content-Type", "application/json") // normal header

		encoder := json.NewEncoder(w)
		if err := encoder.Encode(response); err != nil {
			log.Printf("failed to encode response: %s", err)
			http.Error(w, "unable to encode response", http.StatusInternalServerError)
		}
	})

	return mux
}

func main() {
	// mux := newGeneratorMux()
}

func generateRandomUsers(n int64) []map[string]interface{} {
	var generated []map[string]interface{}

	for i := range n {
		generated = append(generated, map[string]interface{}{
			"username": fmt.Sprintf("user-%d", i+1),
		})
	}

	return generated
}
