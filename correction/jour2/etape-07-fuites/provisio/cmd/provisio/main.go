// provisio : orchestrateur d'activation de services.
//
// Étape 7 : un port admin expose /debug/pprof et un /healthz qui compte
// les goroutines. Le système legacy fictif se tait une fois sur trois : sans
// APP_LEGACY_FIX=1, chaque commande upgrade touchée laisse une goroutine
// bloquée pour toujours.
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
		app.AdminModule,
	).Run()
}
