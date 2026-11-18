// provisio : orchestrateur d'activation de services.
//
// Étape 9 : arrêt propre. SIGTERM → fx arrête les hooks dans l'ordre
// inverse : admin, HTTP (plus de nouvelles commandes), puis le pool qui
// draine sa file et termine les commandes en cours, dans la limite de
// fx.StopTimeout. Au-delà, elles sont annulées avec une cause explicite.
package main

import (
	"log/slog"
	"os"

	"go.uber.org/fx"

	"gitlab.example.com/free/provisio/internal/app"
	"gitlab.example.com/free/provisio/internal/config"
)

func main() {
	// La config est lue avant fx.New : fx.StopTimeout est une option de
	// l'application, pas une dépendance du graphe. fx.Supply la réinjecte.
	cfg, err := config.Load()
	if err != nil {
		slog.Error("configuration invalide", "err", err)
		os.Exit(1)
	}

	fx.New(
		fx.WithLogger(app.FxLogger),
		// Délai total accordé aux OnStop. Au-delà, fx abandonne : les
		// commandes en cours sont annulées par le pool (ErrShuttingDown).
		fx.StopTimeout(cfg.StopTimeout),
		fx.Supply(cfg),

		app.LoggingModule,
		app.ClientsModule,
		app.OrchestratorModule,
		app.WorkerModule,
		app.HTTPModule,
		app.AdminModule,
	).Run()
}
