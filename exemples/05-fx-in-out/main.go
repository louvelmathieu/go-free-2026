// « Six paramètres dans un constructeur, c'est normal ? »
//
// Tant que ça tient, oui : une signature ordinaire est la forme normale.
// fx.In sert quand ça ne tient plus, ou quand il faut un tag sur un paramètre.
//
//	go run ./05-fx-in-out
package main

import (
	"context"
	"fmt"

	"go.uber.org/fx"
)

type Config struct{ Timeout int }
type Inventory struct{}
type Configurator struct{}
type Stock struct{}
type Legacy struct{}
type Billing struct{}
type Metrics struct{}

//  1. La signature ordinaire. Six paramètres : lisible encore, mais chaque
//     nouvelle dépendance touche la signature et tous les appels dans les tests.
func NewOrchestratorV1(cfg Config, inv *Inventory, conf *Configurator, stock *Stock, legacy *Legacy, billing *Billing) string {
	return "v1 : six paramètres, tout est là"
}

// 2. Une struct ordinaire ? fx ne sait pas la remplir : il lui faudrait un
//    constructeur pour Deps, qui aurait lui-même les six paramètres. On a
//    déplacé le problème, pas résolu.

//  3. fx.In : une struct que fx reconnaît et remplit champ par champ depuis le
//     graphe. Ajouter une dépendance = ajouter un champ. Et les tags :
type OrchestratorDeps struct {
	fx.In
	Cfg     Config
	Inv     *Inventory
	Conf    *Configurator
	Stock   *Stock
	Legacy  *Legacy
	Billing *Billing
	Metrics *Metrics `optional:"true"` // peut manquer : reste nil, pas d'erreur
}

func NewOrchestratorV2(d OrchestratorDeps) string {
	return fmt.Sprintf("v2 : fx.In, metrics présent = %v", d.Metrics != nil)
}

//  4. fx.Out, le symétrique : un constructeur rend une struct, chaque champ
//     devient une valeur fournie au graphe. Pour une initialisation qui produit
//     plusieurs choses d'un coup.
type Clients struct {
	fx.Out
	Inv     *Inventory
	Conf    *Configurator
	Stock   *Stock
	Legacy  *Legacy
	Billing *Billing
}

func NewClients(cfg Config) Clients {
	return Clients{Inv: &Inventory{}, Conf: &Configurator{}, Stock: &Stock{}, Legacy: &Legacy{}, Billing: &Billing{}}
}

func main() {
	app := fx.New(
		fx.NopLogger,
		fx.Supply(Config{Timeout: 5}),
		fx.Provide(NewClients), // fournit cinq types en une fois (fx.Out)
		fx.Provide(
			fx.Annotate(NewOrchestratorV1, fx.ResultTags(`name:"v1"`)),
			fx.Annotate(NewOrchestratorV2, fx.ResultTags(`name:"v2"`)),
		),
		fx.Invoke(fx.Annotate(func(v1, v2 string) {
			fmt.Println(v1)
			fmt.Println(v2)
		}, fx.ParamTags(`name:"v1"`, `name:"v2"`))),
	)
	ctx := context.Background()
	if err := app.Start(ctx); err != nil {
		panic(err)
	}
	app.Stop(ctx)
	fmt.Println()
	fmt.Println("Règle : signature ordinaire jusqu'à quatre ou cinq paramètres,")
	fmt.Println("fx.In au-delà, ou dès qu'un paramètre a besoin d'un tag (optional, name, group).")
}
