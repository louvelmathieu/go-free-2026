// Piège 4 : "missing type" dans un fx.Module.
//
// Lancer : go run ./pieges/04-missing-type-module
//
// Le module "testmodule" fournit un Repository dont le constructeur réclame
// un *Config. Personne ne fournit *Config : fx refuse de démarrer avec
// "missing type: *main.Config".
//
// Ici on garde le logger par défaut (pas de fx.NopLogger) pour lire les logs :
//   - chaque PROVIDE indique le module d'origine : (from module "testmodule")
//   - l'erreur donne la chaîne complète : qui demande quoi, dans quel module,
//     et à quelle ligne le constructeur fautif est déclaré.
//
// Correction : fournir *Config, dans le module ou à la racine
// (un module voit ce que fournit son parent, l'inverse n'est vrai que si
// le Provide du module n'est pas marqué fx.Private).
package main

import (
	"fmt"

	"go.uber.org/fx"
)

type Config struct{ DSN string }

type Repository struct{ cfg *Config }

func NewRepository(cfg *Config) *Repository { return &Repository{cfg: cfg} }

type Service struct{ repo *Repository }

func NewService(r *Repository) *Service { return &Service{repo: r} }

var TestModule = fx.Module("testmodule",
	fx.Provide(
		NewRepository,
		NewService,
	),
	// Correction : décommenter la ligne suivante.
	// fx.Provide(func() *Config { return &Config{DSN: "postgres://localhost/provisio"} }),
)

func main() {
	app := fx.New(
		TestModule,
		fx.Invoke(func(s *Service) { fmt.Println("service prêt, DSN =", s.repo.cfg.DSN) }),
	)

	if err := app.Err(); err != nil {
		fmt.Println("\nfx refuse de démarrer :")
		fmt.Println(err)
	}
}
