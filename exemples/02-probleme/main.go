// Le problème : des erreurs qui ne portent rien. Le code "naïf" de la slide,
// qui marche, jusqu'au jour où quelqu'un change un message.
//
//	go run ./02-probleme
package main

import (
	"errors"
	"fmt"
	"strings"
)

// L'inventaire renvoie une erreur avec un message. Juste un message.
func checkEligibility(address string) error {
	if strings.Contains(address, "hors-zone") {
		return errors.New("adresse non éligible")
	}
	return nil
}

// L'orchestrateur la remonte telle quelle : d'où vient-elle ? quelle étape ?
func run(address string) error {
	if err := checkEligibility(address); err != nil {
		return err
	}
	return nil
}

// Le handler HTTP, six mois plus tard, décide du code HTTP... en lisant le texte.
func handle(address string) (int, string) {
	err := run(address)
	if err == nil {
		return 200, "ok"
	}
	if strings.Contains(err.Error(), "non éligible") { // aïe
		return 422, err.Error()
	}
	return 500, err.Error() // et le message interne part au client
}

func main() {
	code, body := handle("hameau hors-zone")
	fmt.Println(code, body)

	// Un jour, quelqu'un traduit le message de l'inventaire en anglais...
	fmt.Println()
	fmt.Println("après un changement de message côté inventaire :")
	fmt.Println(handleAfterRename("hameau hors-zone"))
	fmt.Println("→ le 422 est devenu un 500, aucun test n'a bougé, la prod le découvre.")
}

func handleAfterRename(address string) (int, string) {
	err := errors.New("address not eligible") // même cause, autre texte
	if strings.Contains(err.Error(), "non éligible") {
		return 422, err.Error()
	}
	return 500, err.Error()
}
