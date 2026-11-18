// Package app assemble provisio avec fx. Chaque fx.Module regroupe les
// constructeurs d'une zone du service ; main se contente de les empiler.
package app

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"net/http/pprof"
	"os"
	"runtime"
	"time"

	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"

	"gitlab.example.com/free/freelib/httpx"
	"gitlab.example.com/free/freelib/logging"
	"gitlab.example.com/free/provisio/internal/api"
	"gitlab.example.com/free/provisio/internal/clients/fake"
	"gitlab.example.com/free/provisio/internal/config"
	"gitlab.example.com/free/provisio/internal/orchestrator"
	"gitlab.example.com/free/provisio/internal/store"
)

// Config : un seul constructeur, fx l'appelle une seule fois, tous ceux
// qui demandent config.Config reçoivent la même valeur.
var ConfigModule = fx.Module("config",
	fx.Provide(config.Load),
)

// Logging : le logger dépend de la config. fx déduit l'ordre tout seul.
var LoggingModule = fx.Module("logging",
	fx.Provide(func(cfg config.Config) *slog.Logger {
		return logging.New(os.Stdout, "provisio", cfg.Env, cfg.LogLevel)
	}),
	// Invoke : exécuté au démarrage, même si personne ne consomme le résultat.
	fx.Invoke(func(l *slog.Logger) { slog.SetDefault(l) }),
)

// Clients : les fakes sont fournis sous leur type d'interface, pas sous
// leur type concret. L'orchestrateur demande un orchestrator.Inventory,
// il ne saura jamais qu'il parle à un *fake.Inventory.
var ClientsModule = fx.Module("clients",
	fx.Provide(
		fx.Annotate(
			func() *fake.Inventory { return &fake.Inventory{Latency: 300 * time.Millisecond} },
			fx.As(new(orchestrator.Inventory)),
		),
		fx.Annotate(
			func() *fake.Configurator { return &fake.Configurator{Latency: 200 * time.Millisecond} },
			fx.As(new(orchestrator.Configurator)),
		),
		fx.Annotate(
			func(cfg config.Config) *fake.Stock { return fake.NewStock(cfg.Stock) },
			fx.As(new(orchestrator.Stock)),
		),
		fx.Annotate(
			func(cfg config.Config) *fake.Legacy {
				return &fake.Legacy{
					Latency:        400 * time.Millisecond,
					FailRate:       cfg.LegacyFailRate,
					HangRate:       cfg.LegacyHangRate,
					RespectContext: cfg.LegacyRespectContext,
				}
			},
			fx.As(new(orchestrator.LegacyClient)),
		),
		fx.Annotate(
			func() *fake.Billing { return &fake.Billing{Latency: 100 * time.Millisecond} },
			fx.As(new(orchestrator.Billing)),
		),
	),
)

// Orchestrator : fx.In regroupe les dépendances dans une struct plutôt
// qu'une signature à six paramètres.
type orchestratorDeps struct {
	fx.In

	Config       config.Config
	Inventory    orchestrator.Inventory
	Configurator orchestrator.Configurator
	Stock        orchestrator.Stock
	Legacy       orchestrator.LegacyClient
	Billing      orchestrator.Billing
}

var OrchestratorModule = fx.Module("orchestrator",
	fx.Provide(
		func(d orchestratorDeps) *orchestrator.Orchestrator {
			return orchestrator.New(
				orchestrator.WithInventory(d.Inventory),
				orchestrator.WithConfigurator(d.Configurator),
				orchestrator.WithStock(d.Stock),
				orchestrator.WithLegacy(d.Legacy),
				orchestrator.WithBilling(d.Billing),
				orchestrator.WithStepTimeout(d.Config.StepTimeout),
			)
		},
		// Le store est fourni sous l'interface que l'API attend.
		fx.Annotate(
			func() *store.Memory { return &store.Memory{} },
			fx.As(new(api.Store)),
		),
		// L'orchestrateur concret satisfait api.Runner : on le re-fournit
		// sous ce type sans le reconstruire.
		func(o *orchestrator.Orchestrator) api.Runner { return o },
	),
)

// HTTP : le serveur est attaché au cycle de vie. OnStart ne bloque pas
// (ListenAndServe part dans une goroutine), OnStop arrête proprement.
var HTTPModule = fx.Module("http",
	fx.Provide(
		func(r api.Runner, st api.Store, cfg config.Config) *api.Handler {
			return api.NewHandler(r, st, cfg.OrderTimeout)
		},
		func(cfg config.Config, h *api.Handler, l *slog.Logger) *http.Server {
			return &http.Server{
				Addr:              cfg.Addr,
				Handler:           httpx.RequestLogger(l)(h.Routes()),
				ReadHeaderTimeout: 5 * time.Second,
			}
		},
	),
	fx.Invoke(startServer),
)

func startServer(lc fx.Lifecycle, srv *http.Server, l *slog.Logger) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			// Listen ici, Serve dans la goroutine : si le port est pris,
			// l'erreur remonte tout de suite et fx refuse de démarrer.
			ln, err := net.Listen("tcp", srv.Addr)
			if err != nil {
				return err
			}
			l.Info("provisio en écoute", "addr", ln.Addr().String())
			go func() {
				if err := srv.Serve(ln); err != nil && !errors.Is(err, http.ErrServerClosed) {
					l.Error("serveur HTTP arrêté", "err", err)
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			l.Info("arrêt du serveur HTTP")
			return srv.Shutdown(ctx)
		},
	})
}

// Admin : un second serveur HTTP, sur une adresse locale, pour l'outillage.
// /debug/pprof n'a rien à faire sur le port public.
var AdminModule = fx.Module("admin",
	fx.Invoke(func(lc fx.Lifecycle, cfg config.Config, l *slog.Logger) {
		mux := http.NewServeMux()
		mux.HandleFunc("GET /healthz", func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintf(w, "ok goroutines=%d\n", runtime.NumGoroutine())
		})
		mux.HandleFunc("/debug/pprof/", pprof.Index)
		mux.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
		mux.HandleFunc("/debug/pprof/profile", pprof.Profile)
		mux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
		mux.HandleFunc("/debug/pprof/trace", pprof.Trace)
		srv := &http.Server{Addr: cfg.AdminAddr, Handler: mux, ReadHeaderTimeout: 5 * time.Second}
		lc.Append(fx.Hook{
			OnStart: func(ctx context.Context) error {
				ln, err := net.Listen("tcp", srv.Addr)
				if err != nil {
					return err
				}
				l.Info("admin en écoute", "addr", ln.Addr().String())
				go func() { _ = srv.Serve(ln) }()
				return nil
			},
			OnStop: func(ctx context.Context) error { return srv.Shutdown(ctx) },
		})
	}),
)

// FxLogger envoie les événements internes de fx (provide, invoke, hooks)
// dans notre slog au lieu du stdout brut. Indispensable pour débugger
// un graphe qui ne se résout pas.
func FxLogger(l *slog.Logger) fxevent.Logger {
	fl := &fxevent.SlogLogger{Logger: l.With("component", "fx")}
	fl.UseLogLevel(slog.LevelDebug) // bavard en dev, silencieux en prod
	return fl
}
