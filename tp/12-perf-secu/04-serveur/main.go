// 04-serveur/main.go
// Exercice 4 : un petit serveur HTTP qui marche, avec trois faiblesses.
//
//	go run ./04-serveur            (terminal 1, il écoute sur :8080)
//
// Dans un second terminal, trois essais :
//
//	A. Un corps énorme :
//	   head -c 50000000 /dev/zero | curl -X POST --data-binary @- localhost:8080/import
//
//	B. Un client très lent :
//	   go run ./04-serveur/client-lent
//
//	C. Une panic en arrière-plan :
//	   curl localhost:8080/notifier        puis        curl localhost:8080/sante
//
// À faire : observer ce qui se passe à chaque essai, puis corriger le serveur.
package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
)

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /sante", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "ok")
	})

	// Importe un fichier envoyé par le client.
	mux.HandleFunc("POST /import", func(w http.ResponseWriter, r *http.Request) {
		data, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "lecture impossible", http.StatusBadRequest)
			return
		}
		fmt.Fprintf(w, "%d octets importés\n", len(data))
	})

	// Répond tout de suite, et envoie une notification en arrière-plan.
	mux.HandleFunc("GET /notifier", func(w http.ResponseWriter, r *http.Request) {
		go envoyerNotification(nil)
		fmt.Fprintln(w, "notification en cours")
	})

	log.Println("écoute sur :8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}

// envoyerNotification compte les envois par destinataire.
func envoyerNotification(compteurs map[string]int) {
	compteurs["equipe-support"]++
	log.Println("notification envoyée")
}
