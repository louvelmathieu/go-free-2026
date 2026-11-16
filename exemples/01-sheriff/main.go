// Deux vues d'un même type avec les groupes de sheriff : la liste renvoie
// l'essentiel, le détail renvoie tout.
//
//	go run ./01-sheriff
package main

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/liip/sheriff/v2"
)

type OrderResponse struct {
	ID        string    `json:"id"         groups:"list,detail"`
	Type      string    `json:"type"       groups:"list,detail"`
	Status    string    `json:"status"     groups:"list,detail"`
	CreatedAt time.Time `json:"created_at" groups:"list,detail"`
	Steps     []string  `json:"steps"      groups:"detail"` // détail seulement
}

// writeView : 1. sheriff ne garde que les champs du groupe, 2. encoding/json écrit.
func writeView(data any, group string) error {
	view, err := sheriff.Marshal(&sheriff.Options{Groups: []string{group}}, data)
	if err != nil {
		return err
	}
	return json.NewEncoder(os.Stdout).Encode(view)
}

func main() {
	o := OrderResponse{
		ID: "ord_1", Type: "activation", Status: "done",
		CreatedAt: time.Date(2026, 11, 13, 10, 0, 0, 0, time.UTC),
		Steps:     []string{"eligibility", "resource:RES-12/3", "configure", "billing"},
	}
	fmt.Print("GET /orders       (list)   : ")
	writeView([]OrderResponse{o}, "list")
	fmt.Print("GET /orders/{id}  (detail) : ")
	writeView(o, "detail")
}
