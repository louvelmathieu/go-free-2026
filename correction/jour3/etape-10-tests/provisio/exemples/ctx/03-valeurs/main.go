// Valeurs dans le contexte : clé typée non exportée, accesseurs, et ce
// qu'il ne faut PAS y mettre.
//
// Lancer : go run ./exemples/ctx/03-valeurs
package main

import (
	"context"
	"fmt"
)

// Bon usage : des métadonnées de requête, en lecture seule, optionnelles.
type traceKey struct{}

func WithTraceID(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, traceKey{}, id)
}

func TraceID(ctx context.Context) string {
	id, _ := ctx.Value(traceKey{}).(string) // "" si absent : pas de panique
	return id
}

// Mauvais usage : une dépendance métier cachée dans le contexte.
// Le compilateur ne voit plus que handle() a besoin d'un store.
type storeKey struct{}

type Store interface{ Get(id string) string }

func handleMauvais(ctx context.Context, id string) string {
	st := ctx.Value(storeKey{}).(Store) // panique si absent, invisible dans la signature
	return st.Get(id)
}

func handleBon(ctx context.Context, st Store, id string) string {
	return st.Get(id) // la dépendance est explicite, injectable, testable
}

type fakeStore struct{}

func (fakeStore) Get(id string) string { return "commande " + id }

func main() {
	ctx := WithTraceID(context.Background(), "trace-42")
	fmt.Println("trace :", TraceID(ctx))
	fmt.Println("absent :", TraceID(context.Background()) == "")

	// Collision impossible : une clé string "trace" d'un autre package
	// n'est pas notre traceKey{}.
	ctx2 := context.WithValue(ctx, "trace", "autre") //nolint:staticcheck
	fmt.Println("toujours la nôtre :", TraceID(ctx2))

	fmt.Println(handleBon(ctx, fakeStore{}, "1"))
	fmt.Println(handleMauvais(context.WithValue(ctx, storeKey{}, fakeStore{}), "2"))
}
