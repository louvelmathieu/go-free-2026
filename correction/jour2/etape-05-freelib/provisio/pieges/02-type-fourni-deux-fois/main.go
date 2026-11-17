// Piège 2 : deux constructeurs pour le même type.
//
// Lancer : go run ./pieges/02-type-fourni-deux-fois
//
// Un service a souvent plusieurs *http.Client (un par système distant) ou
// plusieurs *sql.DB. Fournir deux fois le même type est une erreur de
// graphe : fx refuse de démarrer. La solution est de les nommer
// (fx.Annotate + fx.ResultTags / fx.ParamTags), pas de créer des types
// wrapper artificiels.
package main

import (
	"fmt"
	"net/http"
	"time"

	"go.uber.org/fx"
)

func newInventoryClient() *http.Client { return &http.Client{Timeout: 2 * time.Second} }
func newLegacyClient() *http.Client    { return &http.Client{Timeout: 10 * time.Second} }

func main() {
	app := fx.New(
		fx.NopLogger,

		// Version cassée : "cannot provide *http.Client ... already provided"
		fx.Provide(newInventoryClient, newLegacyClient),

		// Correction : nommer chaque instance.
		// fx.Provide(
		//     fx.Annotate(newInventoryClient, fx.ResultTags(`name:"inventory"`)),
		//     fx.Annotate(newLegacyClient, fx.ResultTags(`name:"legacy"`)),
		// ),

		fx.Invoke(
			// et déclarer lequel on veut :
			// fx.Annotate(func(c *http.Client) { ... }, fx.ParamTags(`name:"legacy"`)),
			func(c *http.Client) { fmt.Println("timeout du client reçu :", c.Timeout) },
		),
	)

	if err := app.Err(); err != nil {
		fmt.Println("fx refuse de démarrer :")
		fmt.Println(err)
	}
}
