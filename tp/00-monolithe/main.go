// provisio, version « ça marche » : tout dans un fichier, tout dans main.
//
// C'est le point de départ du TP 1. Le service fonctionne : il reçoit des
// commandes d'activation ou d'upgrade, enchaîne les appels aux systèmes
// externes (simulés), et répond. Il a juste grandi dans un seul package.
//
//	go run .
//	curl -X POST localhost:8080/orders -d '{"type":"activation","customer_id":"C42","target":"12 rue des Lilas"}'
//	curl localhost:8080/orders/<id>
package main

import (
	crand "crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"math/rand/v2"
	"net/http"
	"strings"
	"sync"
	"time"
)

// Tout est global, tout est exporté, tout est modifiable depuis n'importe où.
var (
	Orders     = map[string]*Order{}
	OrdersMu   sync.Mutex
	StockCount = 100
	StockMu    sync.Mutex
	Timeout    = 2 * time.Second
)

type Order struct {
	ID         string    `json:"id"`
	Type       string    `json:"type"`
	CustomerID string    `json:"customer_id"`
	Target     string    `json:"target"`
	Status     string    `json:"status"`
	Steps      []string  `json:"steps"`
	CreatedAt  time.Time `json:"created_at"`
}

func main() {
	http.HandleFunc("POST /orders", createOrder)
	http.HandleFunc("GET /orders/{id}", getOrder)
	http.HandleFunc("GET /healthz", func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(204) })
	log.Println("provisio en écoute sur :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func createOrder(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Type       string `json:"type"`
		CustomerID string `json:"customer_id"`
		Target     string `json:"target"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "corps JSON invalide", 400)
		return
	}
	if req.Type != "activation" && req.Type != "upgrade" {
		http.Error(w, "type inconnu", 400)
		return
	}

	b := make([]byte, 8)
	crand.Read(b)
	o := &Order{"ord_" + hex.EncodeToString(b), req.Type, req.CustomerID, req.Target, "running", nil, time.Now()}
	OrdersMu.Lock()
	Orders[o.ID] = o
	OrdersMu.Unlock()

	// L'orchestration, inline dans le handler.
	var err error
	if o.Type == "activation" {
		// éligibilité
		time.Sleep(300 * time.Millisecond)
		if strings.Contains(o.Target, "hors-zone") {
			err = fmt.Errorf("adresse %s non éligible", o.Target)
		} else {
			o.Steps = append(o.Steps, "eligibility")
			// réservation de ressource dans l'inventaire
			time.Sleep(300 * time.Millisecond)
			res := fmt.Sprintf("RES-%02d/%d", rand.IntN(50), rand.IntN(8))
			o.Steps = append(o.Steps, "resource:"+res)
			// configuration de l'équipement
			time.Sleep(200 * time.Millisecond)
			o.Steps = append(o.Steps, "configure")
			// facturation
			time.Sleep(100 * time.Millisecond)
			o.Steps = append(o.Steps, "billing")
		}
	} else {
		// réservation d'un équipement en stock
		StockMu.Lock()
		if StockCount == 0 {
			err = fmt.Errorf("stock épuisé")
		} else {
			StockCount--
		}
		StockMu.Unlock()
		if err == nil {
			ref := fmt.Sprintf("EQ-%08d", rand.IntN(1_000_000_000))
			o.Steps = append(o.Steps, "equipment:"+ref)
			// provisioning dans le système legacy, qui échoue parfois
			time.Sleep(400 * time.Millisecond)
			if rand.Float64() < 0.1 {
				err = fmt.Errorf("legacy: timeout interne sur %s", o.Target)
				StockMu.Lock()
				StockCount++ // on rend l'équipement... mais seulement ici
				StockMu.Unlock()
			} else {
				o.Steps = append(o.Steps, "legacy")
				time.Sleep(100 * time.Millisecond)
				o.Steps = append(o.Steps, "billing")
			}
		}
	}

	if err != nil {
		o.Status = "failed"
		http.Error(w, err.Error(), 500)
		return
	}
	o.Status = "done"
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(201)
	json.NewEncoder(w).Encode(o)
}

func getOrder(w http.ResponseWriter, r *http.Request) {
	OrdersMu.Lock()
	o, ok := Orders[r.PathValue("id")]
	OrdersMu.Unlock()
	if !ok {
		http.Error(w, "commande introuvable", 404)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(o)
}
