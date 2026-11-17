// Provide est paresseux, Invoke est actif. La différence, exécutée.
//
//	go run ./05-fx-provide-invoke
package main

import (
	"context"
	"fmt"

	"go.uber.org/fx"
)

type Cache struct{}
type Mailer struct{}

func newCache() *Cache {
	fmt.Println("newCache appelé")
	return &Cache{}
}

func newMailer() *Mailer {
	fmt.Println("newMailer appelé") // ne s'affichera jamais : personne ne demande de *Mailer
	return &Mailer{}
}

func main() {
	app := fx.New(
		fx.NopLogger,
		fx.Provide(newCache, newMailer), // deux constructeurs enregistrés, zéro exécuté
		fx.Invoke(func(c *Cache) { // seul *Cache est demandé
			fmt.Println("invoke : j'ai reçu", c)
		}),
	)
	ctx := context.Background()
	if err := app.Start(ctx); err != nil {
		panic(err)
	}
	app.Stop(ctx)

	fmt.Println()
	fmt.Println("newMailer n'a pas été appelé : un Provide sans consommateur ne s'exécute jamais.")
	fmt.Println("Pas d'erreur, pas de log. C'est le piège n°1 de fx.")
}
