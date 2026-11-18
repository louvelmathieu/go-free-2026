// 04-serveur/main.go
// Exercice 4, corrigé. Chaque modification porte le commentaire CORRECTION.
//
//	go run ./04-serveur
package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /sante", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "ok")
	})

	// Importe un fichier envoyé par le client.
	mux.HandleFunc("POST /import", func(w http.ResponseWriter, r *http.Request) {
		// CORRECTION A : on borne la taille du corps à 1 Mo. Au-delà, la
		// lecture s'arrête avec une erreur, sans avoir tout chargé en mémoire.
		r.Body = http.MaxBytesReader(w, r.Body, 1<<20)

		data, err := io.ReadAll(r.Body)
		if err != nil {
			var tropGros *http.MaxBytesError
			if errors.As(err, &tropGros) {
				http.Error(w, "fichier trop volumineux (1 Mo maximum)", http.StatusRequestEntityTooLarge)
				return
			}
			http.Error(w, "lecture impossible", http.StatusBadRequest)
			return
		}
		fmt.Fprintf(w, "%d octets importés\n", len(data))
	})

	// Répond tout de suite, et envoie une notification en arrière-plan.
	mux.HandleFunc("GET /notifier", func(w http.ResponseWriter, r *http.Request) {
		// CORRECTION C (1/2) : la map est créée.
		go envoyerNotification(map[string]int{})
		fmt.Fprintln(w, "notification en cours")
	})

	// CORRECTION B : un serveur avec des délais, à la place de
	// http.ListenAndServe qui n'en a aucun.
	srv := &http.Server{
		Addr:              ":8080",
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,  // le client a 5 s pour envoyer ses en-têtes
		ReadTimeout:       15 * time.Second, // et 15 s pour envoyer toute sa requête
		WriteTimeout:      30 * time.Second, // on a 30 s pour écrire la réponse
		IdleTimeout:       60 * time.Second, // une connexion inactive est fermée au bout d'une minute
	}
	log.Println("écoute sur :8080")
	log.Fatal(srv.ListenAndServe())
}

// envoyerNotification compte les envois par destinataire.
func envoyerNotification(compteurs map[string]int) {
	// CORRECTION C (2/2) : une panic dans une goroutine tue tout le
	// programme. Toute goroutine qu'on lance soi-même doit la rattraper.
	defer func() {
		if r := recover(); r != nil {
			log.Println("panic dans envoyerNotification :", r)
		}
	}()

	compteurs["equipe-support"]++
	log.Println("notification envoyée")
}
