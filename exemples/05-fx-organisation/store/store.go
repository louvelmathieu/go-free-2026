// Package store : chaque package expose ce qu'il apporte au graphe fx,
// et rien d'autre. main n'a pas à connaître ses constructeurs.
package store

import "go.uber.org/fx"

// Store est l'interface que les autres packages consomment.
type Store interface{ Get(id string) string }

type Memory struct{}

func NewMemory() *Memory             { return &Memory{} }
func (*Memory) Get(id string) string { return "commande " + id + " (mémoire)" }

// Module est la seule chose que main importe de ce package. fx.Module rend
// le groupe nommé : le nom "store" apparaît dans les logs et erreurs de fx.
// On fournit le type concret SOUS son interface, avec Annotate + As.
func Module() fx.Option {
	return fx.Module("store",
		fx.Provide(fx.Annotate(NewMemory, fx.As(new(Store)))),
	)
}
