// fx en quelques lignes : on décrit comment fabriquer les pièces (Provide) et
// ce qu'on veut faire avec (Invoke). fx trouve l'ordre tout seul.
//
//	go run ./05-fx-hello
package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"os"

	"go.uber.org/fx"
)

// Un constructeur : ses paramètres sont ce dont il a besoin, son résultat
// est ce qu'il fournit. Ici : rien en entrée, un *slog.Logger en sortie.
func newLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(os.Stdout, nil))
}

// Greeter a besoin d'un logger. Il ne sait pas d'où il vient, et c'est le but.
type Greeter struct{ log *slog.Logger }

func newGreeter(log *slog.Logger) *Greeter {
	log.Info("greeter construit")
	return &Greeter{log: log}
}

func (g *Greeter) Hello(name string) { fmt.Println("bonjour", name) }

func main() {
	app := fx.New(
		fx.NopLogger, // on cache les logs internes de fx pour l'exemple

		// Provide : "voici comment fabriquer un *slog.Logger et un *Greeter".
		// Rien n'est appelé pour l'instant.
		fx.Provide(newLogger, newGreeter),

		// Invoke : "au démarrage, exécute ceci". fx voit qu'il faut un *Greeter,
		// donc un *slog.Logger : il appelle newLogger, puis newGreeter, puis nous.
		fx.Invoke(func(g *Greeter) { g.Hello("Free") }),
	)

	// Dans un vrai service on écrit app.Run(), qui démarre puis attend Ctrl-C.
	// Ici rien ne tourne en fond : on démarre et on arrête tout de suite.
	ctx := context.Background()
	if err := app.Start(ctx); err != nil {
		log.Fatal(err)
	}
	app.Stop(ctx)
}
