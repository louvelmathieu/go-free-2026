// Trois générations de logging, et comment les faire cohabiter.
//
//	go run ./03-log-logrus-slog
package main

import (
	stdlog "log"
	"log/slog"
	"os"

	log "github.com/sirupsen/logrus" // l'astuce historique : logrus reprend l'API de log
)

func main() {
	// 1. logrus importé sous le nom "log" : le code écrit pour le package
	//    standard compile tel quel, et gagne les niveaux et les champs.
	log.Printf("commande %s démarrée (logrus via l'alias log)", "ord_42")
	log.WithFields(log.Fields{"order_id": "ord_42", "step": "legacy"}).Warn("étape lente")

	// 2. slog n'a pas Printf : l'astuce de l'alias ne marche pas. Mais slog.SetDefault
	//    redirige le package log standard vers le handler slog.
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, nil)))
	stdlog.Printf("commande %s démarrée (ancien code log, inchangé)", "ord_42")
	// → sort en JSON, niveau INFO, sans toucher à cette ligne.

	// 3. Le nouveau code écrit directement en slog, avec des attributs.
	slog.Info("commande démarrée", "order_id", "ord_42", "type", "activation")
}
