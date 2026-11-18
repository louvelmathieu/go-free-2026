// provisio : orchestrateur d'activation de services.
//
// Étape 8 : un pool de workers de taille fixe devant une queue bornée,
// des sémaphores par système externe (inventaire réseau : 3, système legacy : 20), load shedding
// en 503 quand la file est pleine, et Order rendu concurrency-safe.
package main

import (
	"go.uber.org/fx"

	"gitlab.example.com/free/provisio/internal/app"
)

func main() {
	fx.New(
		fx.WithLogger(app.FxLogger),

		app.ConfigModule,
		app.LoggingModule,
		app.ClientsModule,
		app.OrchestratorModule,
		app.WorkerModule,
		app.HTTPModule,
		app.AdminModule,
	).Run()
}
