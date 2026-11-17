// fx.Supply : donner au graphe une valeur qu'on a déjà, plutôt qu'un
// constructeur. Trois cas concrets, puis Decorate, Replace et Populate.
//
//	go run ./05-fx-supply-options
package main

import (
	"context"
	"flag"
	"fmt"
	"time"

	"go.uber.org/fx"
)

type Config struct {
	Env     string
	Timeout time.Duration
}

type Store interface{ Get(id string) string }
type Memory struct{ prefix string }

func (m *Memory) Get(id string) string { return m.prefix + id }

type Handler struct {
	store Store
	cfg   Config
}

func NewHandler(s Store, cfg Config) *Handler { return &Handler{store: s, cfg: cfg} }

func main() {
	env := flag.String("env", "dev", "environnement")
	flag.Parse()

	// Cas 1 : une valeur construite AVANT fx.New, parce qu'on en a besoin pour
	// régler fx lui-même (StopTimeout) ou parce qu'elle vient d'ailleurs (flags).
	cfg := Config{Env: *env, Timeout: 2 * time.Second}

	// Cas 2 : une valeur déjà construite qu'on veut fournir SOUS UNE INTERFACE.
	// Supply(mem) fournirait un *Memory, que personne ne demande. On annote :
	mem := &Memory{prefix: "commande "}

	var handler *Handler

	app := fx.New(
		fx.NopLogger,
		fx.Supply(cfg), // fournit Config
		fx.Supply(fx.Annotate(mem, fx.As(new(Store)))), // fournit Store, pas *Memory
		// Cas 3, la forme avec un nom, si deux valeurs du même type coexistent :
		// fx.Supply(fx.Annotated{Name: "replica", Target: mem}),

		fx.Provide(NewHandler),

		// Decorate : modifier une valeur du graphe pour ceux qui viennent après.
		fx.Decorate(func(cfg Config) Config {
			if cfg.Env == "dev" {
				cfg.Timeout *= 2
			}
			return cfg
		}),

		// Populate : récupérer une valeur du graphe dans une variable à nous.
		fx.Populate(&handler),

		fx.StopTimeout(cfg.Timeout),
	)
	ctx := context.Background()
	if err := app.Start(ctx); err != nil {
		panic(err)
	}
	app.Stop(ctx)
	fmt.Println("handler.store.Get :", handler.store.Get("ord_42"))
	fmt.Println("handler.cfg.Timeout :", handler.cfg.Timeout, "(2s doublé par Decorate en dev)")

	// Replace : substituer une valeur fournie, sans toucher au module qui la
	// fournit. L'outil des tests (module 10).
	var h2 *Handler
	fx.New(fx.NopLogger,
		fx.Supply(cfg), fx.Supply(fx.Annotate(mem, fx.As(new(Store)))), fx.Provide(NewHandler),
		fx.Replace(fx.Annotate(&Memory{prefix: "FAKE "}, fx.As(new(Store)))),
		fx.Populate(&h2),
	).Start(ctx)
	fmt.Println("avec Replace        :", h2.store.Get("ord_42"))
}
