// provisio : orchestrateur d'activation de services.
//
// Étape 1 : tout est câblé à la main dans main. C'est volontaire, on verra
// en fin de journée ce que fx apporte par rapport à ça.
package main

import (
	"log"
	"net/http"
	"time"

	"gitlab.example.com/free/provisio/internal/api"
	"gitlab.example.com/free/provisio/internal/clients/fake"
	"gitlab.example.com/free/provisio/internal/orchestrator"
	"gitlab.example.com/free/provisio/internal/store"
)

func main() {
	orch := orchestrator.New(
		orchestrator.WithInventory(&fake.Inventory{Latency: 300 * time.Millisecond}),
		orchestrator.WithConfigurator(&fake.Configurator{Latency: 200 * time.Millisecond}),
		orchestrator.WithStock(fake.NewStock(2)),
		orchestrator.WithLegacy(&fake.Legacy{Latency: 400 * time.Millisecond, FailRate: 0.1}),
		orchestrator.WithBilling(&fake.Billing{Latency: 100 * time.Millisecond}),
		orchestrator.WithStepTimeout(2*time.Second),
	)

	var orders store.Memory // valeur vierge utilisable
	handler := api.NewHandler(orch, &orders)

	srv := &http.Server{
		Addr:    ":8080",
		Handler: handler.Routes(),
	}

	log.Println("provisio en écoute sur :8080")
	log.Fatal(srv.ListenAndServe())
}
