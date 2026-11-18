// 01-bench/parse.go
// Exercice 1 : deux façons de décoder une commande. Laquelle est la plus
// rapide, de combien, et pourquoi ?
//
// À faire :
//  1. écrire parse_test.go avec un benchmark par fonction,
//  2. lancer : go test ./01-bench -bench . -benchmem
//  3. lire ns/op, B/op et allocs/op,
//  4. écrire ParseCSVNoAlloc, une version de ParseCSV qui n'alloue rien,
//     et mesurer le gain.
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
func ParseCSV(s string) (CreateOrderRequest, error) {
	parts := strings.Split(s, ";")
	if len(parts) != 3 {
		return CreateOrderRequest{}, errors.New("3 champs attendus")
	}
	return CreateOrderRequest{Type: parts[0], CustomerID: parts[1], Target: parts[2]}, nil
}

// Les données à utiliser dans vos benchmarks.
var (
	JSONIn = []byte(`{"type":"activation","customer_id":"C42","target":"12 rue des Lilas, Lyon"}`)
	CSVIn  = "activation;C42;12 rue des Lilas, Lyon"
)
