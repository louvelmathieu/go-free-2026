// provisio : orchestrateur d'activation de services.
//
// Étape 3 : log/slog remplace le package log. Le logger est construit ici,
// injecté dans le middleware HTTP, et voyage ensuite dans le contexte.
package main

import (
	"log/slog"
	"net/http"
	"os"
	"time"

	"gitlab.example.com/free/provisio/internal/api"
	"gitlab.example.com/free/provisio/internal/clients/fake"
	"gitlab.example.com/free/provisio/internal/logging"
	"gitlab.example.com/free/provisio/internal/orchestrator"
	"gitlab.example.com/free/provisio/internal/store"
)

func main() {
	env := os.Getenv("APP_ENV") // "" ou "dev" : texte, "prod" : JSON
	logger := logging.New(os.Stdout, env, slog.LevelDebug)
	slog.SetDefault(logger) // pour le code qui n'a pas (encore) de contexte

	orch := orchestrator.New(
		orchestrator.WithInventory(&fake.Inventory{Latency: 300 * time.Millisecond}),
		orchestrator.WithConfigurator(&fake.Configurator{Latency: 200 * time.Millisecond}),
		orchestrator.WithStock(fake.NewStock(100)),
		orchestrator.WithLegacy(&fake.Legacy{Latency: 400 * time.Millisecond, FailRate: 0.1}),
		orchestrator.WithBilling(&fake.Billing{Latency: 100 * time.Millisecond}),
		orchestrator.WithStepTimeout(2*time.Second),
	)

	var orders store.Memory
	handler := api.NewHandler(orch, &orders)

	srv := &http.Server{
		Addr:    ":8080",
		Handler: api.RequestLogger(logger)(handler.Routes()),
	}

	logger.Info("provisio en écoute", "addr", srv.Addr, "env", env)
	if err := srv.ListenAndServe(); err != nil {
		logger.Error("serveur arrêté", "err", err)
		os.Exit(1)
	}
}
