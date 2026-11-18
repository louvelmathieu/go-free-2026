// 01-bench/parse.go
// Exercice 1, corrigé : les benchmarks sont dans parse_test.go, et
// ParseCSVNoAlloc est la version sans allocation.
//
//	go test ./01-bench -bench . -benchmem
package bench

import (
	"encoding/json"
	"errors"
	"strings"
)

type CreateOrderRequest struct {
	Type       string `json:"type"`
	CustomerID string `json:"customer_id"`
	Target     string `json:"target"`
}

// ParseJSON : la voie normale.
func ParseJSON(b []byte) (CreateOrderRequest, error) {
	var r CreateOrderRequest
	err := json.Unmarshal(b, &r)
	return r, err
}

// ParseCSV : un format "type;client;cible".
// strings.Split crée un slice pour ranger les trois morceaux : 1 allocation.
func ParseCSV(s string) (CreateOrderRequest, error) {
	parts := strings.Split(s, ";")
	if len(parts) != 3 {
		return CreateOrderRequest{}, errors.New("3 champs attendus")
	}
	return CreateOrderRequest{Type: parts[0], CustomerID: parts[1], Target: parts[2]}, nil
}

// CORRECTION : ParseCSVNoAlloc fait le même travail sans Split. On cherche
// la position des deux points-virgules, puis on découpe la string avec des
// tranches s[a:b], qui ne copient rien : 0 allocation.
func ParseCSVNoAlloc(s string) (CreateOrderRequest, error) {
	i := strings.IndexByte(s, ';')
	if i < 0 {
		return CreateOrderRequest{}, errors.New("3 champs attendus")
	}
	j := strings.IndexByte(s[i+1:], ';')
	if j < 0 {
		return CreateOrderRequest{}, errors.New("3 champs attendus")
	}
	j += i + 1
	return CreateOrderRequest{Type: s[:i], CustomerID: s[i+1 : j], Target: s[j+1:]}, nil
}

// Les données utilisées dans les benchmarks.
var (
	JSONIn = []byte(`{"type":"activation","customer_id":"C42","target":"12 rue des Lilas, Lyon"}`)
	CSVIn  = "activation;C42;12 rue des Lilas, Lyon"
)
