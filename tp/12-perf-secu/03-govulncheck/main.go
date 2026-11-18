// 03-govulncheck/main.go
// Exercice 3 : ce programme marche. Il choisit la langue de la réponse
// d'après l'en-tête Accept-Language envoyé par le navigateur.
//
//	go run ./03-govulncheck
//
// À faire :
//  1. installer l'outil : go install golang.org/x/vuln/cmd/govulncheck@latest
//  2. lancer : govulncheck ./03-govulncheck
//  3. lire le rapport : quelle faille, dans quel module, corrigée dans
//     quelle version, et par quelle ligne de CE fichier on y arrive,
//  4. corriger, puis relancer govulncheck.
package main

import (
	"fmt"

	"golang.org/x/text/language"
)

// LangueDuClient renvoie la langue préférée annoncée par le navigateur.
// L'en-tête vient du client : c'est une donnée qu'on ne contrôle pas.
func LangueDuClient(acceptLanguage string) string {
	tags, _, err := language.ParseAcceptLanguage(acceptLanguage)
	if err != nil || len(tags) == 0 {
		return "fr"
	}
	base, _ := tags[0].Base()
	return base.String()
}

func main() {
	fmt.Println(LangueDuClient("en-GB,en;q=0.9,fr;q=0.8"))
	fmt.Println(LangueDuClient("fr-FR,fr;q=0.9"))
	fmt.Println(LangueDuClient(""))
}
