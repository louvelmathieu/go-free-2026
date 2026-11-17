// Consommer ce qui est dans le graphe : les quatre façons de recevoir une
// dépendance, une par bloc. Il n'y a jamais de container.Get(...) : on
// déclare ce dont on a besoin, et c'est fx qui le pousse.
//
//	go run ./05-fx-consommer
package main

import (
	"context"
	"fmt"

	"go.uber.org/fx"
)

type Store interface{ Get(id string) string }

type Memory struct{ name string }

func (m *Memory) Get(id string) string { return m.name + ":" + id }

// newStore fournit un Store (le type concret *Memory, vu comme Store).
func newStore() Store { return &Memory{name: "memoire"} }

// run démarre et arrête tout de suite une application fx : pas de serveur ici.
func run(opts ...fx.Option) {
	app := fx.New(append([]fx.Option{fx.NopLogger}, opts...)...)
	ctx := context.Background()
	if err := app.Start(ctx); err != nil {
		panic(err)
	}
	app.Stop(ctx)
}

// ---------------------------------------------------------------------------
// 1. Les paramètres d'un constructeur. La façon normale, 90 % du code.
//    Le Handler déclare "j'ai besoin d'un Store" dans sa signature. Il ne sait
//    pas lequel, ni d'où il vient.
// ---------------------------------------------------------------------------

type Handler struct{ store Store }

func NewHandler(s Store) *Handler { return &Handler{store: s} }

func facon1() {
	run(
		fx.Provide(newStore, NewHandler),
		fx.Invoke(func(h *Handler) {
			fmt.Println("1. constructeur :", h.store.Get("ord_42"))
		}),
	)
}

// ---------------------------------------------------------------------------
// 2. Les paramètres d'un Invoke. Même mécanisme, pour une action au démarrage.
//    Ici l'Invoke demande directement le Store, sans passer par un Handler.
// ---------------------------------------------------------------------------

func facon2() {
	run(
		fx.Provide(newStore),
		fx.Invoke(func(s Store) {
			fmt.Println("2. invoke       :", s.Get("ord_42"))
		}),
	)
}

// ---------------------------------------------------------------------------
// 3. Une struct fx.In. Quand la liste de paramètres s'allonge, ou qu'il faut
//    un tag : optional pour une dépendance qui peut manquer.
// ---------------------------------------------------------------------------

type Metrics struct{}

type ReportDeps struct {
	fx.In
	Store   Store
	Metrics *Metrics `optional:"true"` // personne ne fournit *Metrics : reste nil, pas d'erreur
}

func NewReport(d ReportDeps) string {
	return fmt.Sprintf("store=%s metrics=%v", d.Store.Get("ord_42"), d.Metrics)
}

func facon3() {
	run(
		fx.Provide(newStore, NewReport),
		fx.Invoke(func(report string) {
			fmt.Println("3. fx.In        :", report)
		}),
	)
}

// ---------------------------------------------------------------------------
// 4. Depuis l'extérieur : fx.Populate. Après construction, fx copie une valeur
//    du graphe dans une variable à nous. Pour un main qui veut garder la main,
//    ou un test. À utiliser rarement.
// ---------------------------------------------------------------------------

func facon4() {
	var h *Handler
	run(
		fx.Provide(newStore, NewHandler),
		fx.Populate(&h),
	)
	// On est sorti de fx : h est une variable ordinaire, remplie par le graphe.
	fmt.Println("4. populate     :", h.store.Get("ord_42"))
}

func main() {
	facon1()
	facon2()
	facon3()
	facon4()
}
