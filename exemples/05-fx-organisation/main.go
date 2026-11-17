// Organiser une application fx : chaque package expose une fonction Module()
// qui rend un fx.Option. main ne fait que les assembler.
//
//	go run ./05-fx-organisation
package main

import (
	"context"

	"go.uber.org/fx"

	"gitlab.example.com/free/exemples/05-fx-organisation/httpapi"
	"gitlab.example.com/free/exemples/05-fx-organisation/store"
)

func main() {
	// fx.Option (singulier) : le type de tout ce qu'on passe à fx.New
	// (un Provide, un Invoke, un Module, un Supply... sont des fx.Option).
	// fx.Options (pluriel) : la fonction qui en regroupe plusieurs en un seul.
	modules := fx.Options(
		store.Module(),
		httpapi.Module(),
	)

	app := fx.New(fx.NopLogger, modules)
	ctx := context.Background()
	if err := app.Start(ctx); err != nil {
		panic(err)
	}
	app.Stop(ctx)
	// Ajouter un package = une ligne ici. Retirer store.Module() : fx dit
	// « missing type store.Store, requested by httpapi.NewHandler ».
}
