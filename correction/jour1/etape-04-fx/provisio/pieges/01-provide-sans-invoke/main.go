// Piège 1 : un Provide que personne ne consomme ne s'exécute jamais.
//
// Lancer : go run ./pieges/01-provide-sans-invoke
//
// On s'attend à voir "serveur démarré". Il ne se passe rien : fx construit
// le graphe paresseusement, uniquement ce dont les Invoke ont besoin.
// Sans Invoke, aucun constructeur n'est appelé, et fx s'arrête sans erreur.
package main

import (
	"context"
	"fmt"

	"go.uber.org/fx"
)

type Server struct{}

func NewServer() *Server {
	fmt.Println("serveur démarré") // jamais affiché
	return &Server{}
}

func main() {
	app := fx.New(
		fx.NopLogger,
		fx.Provide(NewServer),
		// Correction : décommenter la ligne suivante.
		// fx.Invoke(func(*Server) {}),
	)

	ctx := context.Background()
	if err := app.Start(ctx); err != nil {
		panic(err)
	}
	_ = app.Stop(ctx)
	fmt.Println("fx a démarré puis s'est arrêté sans jamais appeler NewServer")
}
