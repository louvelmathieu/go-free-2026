// 03-pieges-slog/main.go
// Trois pièges slog à voir de ses yeux : clés impaires (!BADKEY),
// With qui fige une valeur, logger nil.
//
//	go run ./03-pieges-slog
//	go vet ./03-pieges-slog     (vet signale le piège 1 avant l'exécution)
package main

import (
	"fmt"
	"log/slog"
	"os"
)

type Order struct {
	ID     string
	Status string
}

func (o *Order) LogValue() slog.Value {
	return slog.GroupValue(slog.String("id", o.ID), slog.String("status", o.Status))
}

func main() {
	// Sans l'heure, pour que la sortie soit lisible à l'écran.
	log := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		ReplaceAttr: func(_ []string, a slog.Attr) slog.Attr {
			if a.Key == slog.TimeKey {
				return slog.Attr{}
			}
			return a
		},
	}))

	fmt.Println("--- Piège 1 : clés impaires ---")
	// slog attend des PAIRES après le message : clé, valeur, clé, valeur...
	log.Info("commande reçue", "order_id", "ord_42") // correct : 1 clé + 1 valeur

	// Ici on a oublié la valeur : "order_id" est seul.
	// slog ne peut pas deviner, il range la string sous la clé !BADKEY.
	log.Info("commande reçue", "order_id")

	// Variante fréquente : une valeur oubliée au milieu décale tout le reste.
	log.Info("commande reçue", "order_id", "customer", "C42")
	//                          clé         valeur      → seul → !BADKEY=C42

	// La parade : les attributs typés. Clé et valeur sont dans le même appel,
	// on ne peut plus en oublier une, et le compilateur vérifie le type.
	log.Info("commande reçue", slog.String("order_id", "ord_42"), slog.Int("attempt", 1))

	fmt.Println("\n--- Piège 2 : With fige la valeur ---")
	o := &Order{ID: "ord_42", Status: "running"}
	fige := log.With("order", o) // LogValue() est appelé ICI, une fois pour toutes
	o.Status = "done"
	fige.Info("via With (statut figé)")
	log.Info("à chaque appel (statut à jour)", "order", o)

	fmt.Println("\n--- Piège 3 : logger nil ---")
	var absent *slog.Logger // un champ de struct jamais initialisé, par exemple
	func() {
		defer func() { fmt.Println("panic rattrapé :", recover()) }()
		absent.Info("jamais affiché")
	}()
}
