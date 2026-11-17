// Et à l'intérieur ? fx n'injecte qu'à la frontière (la signature). Dedans,
// c'est du Go ordinaire : une fonction reçoit ce dont elle a besoin de son
// appelant, ou le lit dans la struct qui le porte.
//
//	go run ./05-fx-interieur
package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"go.uber.org/fx"
)

type Metrics struct{ count int }

func (m *Metrics) Inc() { m.count++ }

type Handler struct{}

// --- Façon 1 : l'Invoke demande tout et transmet à la main ---
// Si registerRoutes a besoin demain d'une chose de plus, l'Invoke doit la
// demander aussi : il connaît les besoins de tout ce qu'il appelle.
func registerRoutes(h *Handler, l *slog.Logger) { l.Info("routes enregistrées (façon 1)") }

// --- Façon 2 : les dépendances vivent dans une struct construite par fx ---
type Server struct {
	handler *Handler
	log     *slog.Logger
	metrics *Metrics
}

// fx remplit ICI, une fois. Ensuite plus personne ne repasse les dépendances.
func NewServer(h *Handler, l *slog.Logger, m *Metrics) *Server {
	return &Server{handler: h, log: l, metrics: m}
}

func (s *Server) RegisterRoutes() {
	s.metrics.Inc() // la dépendance est un champ : rien à passer
	s.log.Info("routes enregistrées (façon 2)", "metrics", s.metrics.count)
}

func main() {
	app := fx.New(
		fx.NopLogger,
		fx.Provide(
			func() *slog.Logger { return slog.New(slog.NewTextHandler(os.Stdout, nil)) },
			func() *Metrics { return &Metrics{} },
			func() *Handler { return &Handler{} },
			NewServer,
		),

		// Façon 1 : tout dans la signature de l'Invoke, puis on transmet.
		fx.Invoke(func(h *Handler, l *slog.Logger) { registerRoutes(h, l) }),

		// Façon 2 : l'Invoke ne connaît que le composant du haut. Les besoins
		// de RegisterRoutes sont réglés dans NewServer, par fx.
		fx.Invoke(func(s *Server) { s.RegisterRoutes() }),
	)
	ctx := context.Background()
	if err := app.Start(ctx); err != nil {
		panic(err)
	}
	app.Stop(ctx)
	fmt.Println("Pas de fx.Get() au fond d'une fonction : si une méthode a besoin d'une")
	fmt.Println("nouvelle dépendance, on l'ajoute au constructeur et aucun appelant ne change.")
}
