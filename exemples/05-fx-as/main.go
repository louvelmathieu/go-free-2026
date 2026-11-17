// « Mon fake est bien dans le graphe, pourquoi fx dit "missing type Inventory" ? »
//
// fx relie les constructeurs par type EXACT. Le constructeur rend un
// *FakeInventory, l'orchestrateur demande une interface Inventory : pour fx
// ce sont deux types sans rapport. Ce programme montre l'erreur, puis les
// deux façons de faire le pont.
//
//	go run ./05-fx-as
package main

import (
	"context"
	"fmt"
	"strings"

	"go.uber.org/fx"
)

// Côté consommateur : l'orchestrateur déclare ce dont il a besoin, une interface.
type Inventory interface{ CheckEligibility(address string) bool }

type Orchestrator struct{ inv Inventory }

func NewOrchestrator(inv Inventory) *Orchestrator { return &Orchestrator{inv: inv} }

// Côté fournisseur : un type concret, dont le constructeur rend le type concret.
// C'est la règle Go : accepter des interfaces, rendre des structs.
type FakeInventory struct{ Latency int }

func NewFakeInventory() *FakeInventory                  { return &FakeInventory{Latency: 300} }
func (f *FakeInventory) CheckEligibility(a string) bool { return !strings.Contains(a, "hors-zone") }

func try(label string, opts ...fx.Option) {
	app := fx.New(append([]fx.Option{fx.NopLogger}, opts...)...)
	if err := app.Start(context.Background()); err != nil {
		// On ne garde que la première ligne du message : c'est elle qui compte.
		fmt.Printf("%-28s ERREUR : %s\n", label, strings.SplitN(err.Error(), "\n", 2)[0])
		return
	}
	app.Stop(context.Background())
	fmt.Printf("%-28s OK\n", label)
}

func main() {
	use := fx.Invoke(func(o *Orchestrator) {})

	// 1. Provide simple : fx connaît un *FakeInventory, personne n'en veut.
	//    L'orchestrateur veut un Inventory : fx ne fait pas le lien tout seul.
	try("1. Provide simple", fx.Provide(NewFakeInventory, NewOrchestrator), use)

	// 2. Un constructeur d'adaptation : une fonction qui prend le concret et
	//    rend l'interface. Ça marche. C'est une fonction de plus par client.
	try("2. constructeur d'adaptation",
		fx.Provide(NewFakeInventory, func(f *FakeInventory) Inventory { return f }, NewOrchestrator), use)

	// 3. fx.Annotate + fx.As : « fournis ce constructeur sous cette interface ».
	//    La même chose que 2, en une ligne, sans fonction à écrire.
	try("3. fx.Annotate + fx.As",
		fx.Provide(fx.Annotate(NewFakeInventory, fx.As(new(Inventory))), NewOrchestrator), use)

	fmt.Println()
	fmt.Println("Pourquoi ne pas faire rendre Inventory à NewFakeInventory directement ?")
	fmt.Println("- le package du fake dépendrait de son consommateur (import dans le mauvais sens)")
	fmt.Println("- on perdrait le type concret : plus d'accès à Latency dans les tests")
	fmt.Println("- un fake qui sert deux consommateurs aux interfaces différentes ne pourrait en satisfaire qu'un")
}
