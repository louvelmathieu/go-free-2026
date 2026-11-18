// Package bench : deux façons de décoder une commande, et un benchmark
// pour les départager avec des chiffres plutôt qu'une opinion.
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

// ParseCSV : un format "type;client;cible" que l'inventaire réseau sait émettre.
// Question : est-ce vraiment plus rapide, et de combien ?
func ParseCSV(s string) (CreateOrderRequest, error) {
	parts := strings.Split(s, ";")
	if len(parts) != 3 {
		return CreateOrderRequest{}, errors.New("3 champs attendus")
	}
	return CreateOrderRequest{Type: parts[0], CustomerID: parts[1], Target: parts[2]}, nil
}

// ParseCSVNoAlloc : la même chose sans Split (zéro allocation).
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
