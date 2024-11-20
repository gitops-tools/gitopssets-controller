package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	gitopssetsv1 "github.com/gitops-tools/gitopssets-controller/api/v1alpha1"
	"github.com/google/go-cmp/cmp"
)

func TestGeneration(t *testing.T) {
	mux := newGeneratorMux()
	server := httptest.NewServer(mux)
	defer server.Close()

	b, err := json.Marshal(gitopssetsv1.GeneratorRequest{
		GeneratorParams: map[string]interface{}{
			"count": 5,
		},
	})
	if err != nil {
		t.Fatal(err)
	}

	req, err := http.NewRequest(http.MethodPost, server.URL+"/generate", bytes.NewReader(b))
	if err != nil {
		t.Fatal(err)
	}

	resp, err := server.Client().Do(req)
	if err != nil {
		t.Fatal(err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("response StatusCode = %v, want %v", resp.StatusCode, http.StatusOK)
	}

	var response gitopssetsv1.GeneratorResponse
	decoder := json.NewDecoder(resp.Body)
	if err := decoder.Decode(&response); err != nil {
		t.Fatal(err)
	}

	want := gitopssetsv1.GeneratorResponse{
		Generated: []map[string]any{
			{"username": string("user-1")}, {"username": string("user-2")},
			{"username": string("user-3")}, {"username": string("user-4")},
			{"username": string("user-5")},
		},
	}
	if diff := cmp.Diff(want, response); diff != "" {
		t.Errorf("failed response:\n%s", diff)
	}
}
