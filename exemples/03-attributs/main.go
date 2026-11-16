// Attributs, With, groupes, LogValuer, et le piège du With figé.
//
//	go run ./03-attributs
package main

import (
	"log/slog"
	"os"
)

type Order struct {
	ID     string
	Status string
}

// LogValue : Order décide lui-même ce qu'il montre dans les logs.
func (o *Order) LogValue() slog.Value {
	return slog.GroupValue(slog.String("id", o.ID), slog.String("status", o.Status))
}

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil)).With("service", "provisio")

	// Attributs par paires, ou typés (plus sûr, pas d'allocation inutile) :
	logger.Info("activation démarrée", "order_id", "ord_42", slog.Int("attempt", 1))

	// Group : des attributs regroupés sous une clé.
	logger.Info("requête", slog.Group("http", slog.String("method", "POST"), slog.Int("status", 202)))

	// With : un logger enfant qui répète ses attributs sur chaque ligne.
	reqLog := logger.With("request_id", "req_7f3a")
	reqLog.Info("début")
	reqLog.Info("fin", "duration_ms", 12)

	// Le piège : With fige la valeur au moment de l'appel.
	o := &Order{ID: "ord_42", Status: "running"}
	frozen := logger.With("order", o) // LogValue est résolu ICI
	o.Status = "done"
	frozen.Info("après Complete, via With      ")             // affiche encore running
	logger.Info("après Complete, à chaque appel", "order", o) // affiche done
}
