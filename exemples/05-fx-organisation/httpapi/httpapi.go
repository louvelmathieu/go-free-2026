package httpapi

import (
	"fmt"

	"go.uber.org/fx"

	"gitlab.example.com/free/exemples/05-fx-organisation/store"
)

type Handler struct{ store store.Store }

func NewHandler(s store.Store) *Handler { return &Handler{store: s} }

func (h *Handler) Get(id string) { fmt.Println("GET /orders/"+id+" →", h.store.Get(id)) }

// Module : ce package fournit un Handler, et déclenche son enregistrement.
// Il dépend d'un store.Store, sans savoir lequel : c'est le graphe qui relie.
func Module() fx.Option {
	return fx.Module("httpapi",
		fx.Provide(NewHandler),
		fx.Invoke(func(h *Handler) { h.Get("ord_42") }),
	)
}
