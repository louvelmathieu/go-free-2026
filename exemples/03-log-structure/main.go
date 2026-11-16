// Un log classique est une phrase, un log structuré est un événement.
//
//	go run ./03-log-structure
package main

import (
	"log"
	"log/slog"
	"os"
	"time"
)

func main() {
	id, step, d := "ord_42", "legacy", 2001*time.Millisecond

	// Avant : une phrase. Pour retrouver tous les échecs à l'étape legacy : une regex.
	log.Printf("commande %s échouée à l'étape %s après %v", id, step, d)

	// Après : un message court et des paires clé / valeur, en texte (dev)...
	text := slog.New(slog.NewTextHandler(os.Stdout, nil))
	text.Error("étape échouée", "order_id", id, "step", step, "duration", d)

	// ...ou en JSON (prod), sans changer l'appel : seul le Handler change.
	js := slog.New(slog.NewJSONHandler(os.Stdout, nil)).With("service", "provisio")
	js.Error("étape échouée", "order_id", id, "step", step, "duration", d)
}
