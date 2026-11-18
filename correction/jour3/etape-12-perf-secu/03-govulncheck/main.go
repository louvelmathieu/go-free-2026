// 03-govulncheck/main.go
// Exercice 3, corrigé. La seule modification est dans go.mod :
// golang.org/x/text est passé à une version corrigée.
//
// Il choisit la langue de la réponse
// d'après l'en-tête Accept-Language envoyé par le navigateur.
//
//	go run ./03-govulncheck
//
//	govulncheck ./03-govulncheck     → No vulnerabilities found.
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
