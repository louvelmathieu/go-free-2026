// Piège 3 : cycle de dépendances.
//
// Lancer : go run ./pieges/03-cycle
//
// L'orchestrateur a besoin du store, le store veut notifier l'orchestrateur
// quand une commande change : A dépend de B qui dépend de A. fx détecte le
// cycle et affiche la chaîne complète. La solution est toujours la même :
// casser le cycle avec une interface, un callback ou un événement, jamais
// avec un "setter" appelé après coup.
package main

import (
	"fmt"

	"go.uber.org/fx"
)

type Store struct{ orch *Orchestrator }
type Orchestrator struct{ store *Store }

func NewStore(o *Orchestrator) *Store        { return &Store{orch: o} }
func NewOrchestrator(s *Store) *Orchestrator { return &Orchestrator{store: s} }

func main() {
	app := fx.New(
		fx.NopLogger,
		fx.Provide(NewStore, NewOrchestrator),
		fx.Invoke(func(*Orchestrator) {}),
	)

	if err := app.Err(); err != nil {
		fmt.Println("fx refuse de démarrer :")
		fmt.Println(err)
	}
}
