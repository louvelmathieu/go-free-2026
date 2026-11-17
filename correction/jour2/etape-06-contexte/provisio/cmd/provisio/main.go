// provisio : orchestrateur d'activation de services.
//
// Étape 5 : fx assemble le service. main ne construit plus rien, il décrit
// quels modules composent l'application et laisse fx résoudre le graphe,
// démarrer les hooks, puis les arrêter dans l'ordre inverse sur SIGTERM.
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
		app.HTTPModule,
	).Run()
}
