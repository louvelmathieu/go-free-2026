// Lifecycle : les OnStart dans l'ordre, les OnStop dans l'ordre inverse.
//
//	go run ./05-fx-lifecycle
package main

import (
	"context"
	"fmt"

	"go.uber.org/fx"
)

type DB struct{}
type Server struct{ db *DB }

func newDB(lc fx.Lifecycle) *DB {
	lc.Append(fx.Hook{
		OnStart: func(context.Context) error { fmt.Println("1. DB : connexion ouverte"); return nil },
		OnStop:  func(context.Context) error { fmt.Println("4. DB : connexion fermée"); return nil },
	})
	return &DB{}
}

func newServer(lc fx.Lifecycle, db *DB) *Server {
	lc.Append(fx.Hook{
		// OnStart doit rendre la main : ce qui bloque part dans une goroutine.
		OnStart: func(context.Context) error { fmt.Println("2. serveur : en écoute"); return nil },
		OnStop:  func(context.Context) error { fmt.Println("3. serveur : arrêté"); return nil },
	})
	return &Server{db: db}
}

func main() {
	app := fx.New(
		fx.NopLogger,
		fx.Provide(newDB, newServer),
		fx.Invoke(func(*Server) {}), // on tire le graphe, rien d'autre
	)
	ctx := context.Background()
	if err := app.Start(ctx); err != nil {
		panic(err)
	}
	fmt.Println("   ... le service tourne (ici : rien à faire) ...")
	app.Stop(ctx)
	// Le serveur dépend de la DB : il démarre après elle, et s'arrête avant elle.
	// C'est exactement l'ordre qu'on veut, sans l'avoir écrit.
}
