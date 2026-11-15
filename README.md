# Formation Go avancé, fil rouge « provisio »

Support de code de la formation Go avancé (Sparks Formation, 3 jours).
Le fil rouge est **provisio**, un orchestrateur d'activation de services :
une API reçoit des commandes, un orchestrateur enchaîne les appels à des
systèmes internes simulés (inventaire, configuration, stock d'équipements,
système legacy, facturation). Le domaine est volontairement simple, la
structure ressemble à un vrai service.

## Organisation

```
tp/00-monolithe/          point de départ du TP 1 : le service dans un seul main.go
exemples/                 petits programmes autonomes montrés dans les slides (go run ./xx)
correction/
  jour1/
    etape-01-structure/   module 1  (et correction du TP 1)
    etape-02-erreurs/     module 2
    etape-03-slog/        module 3
    etape-04-fx/          module 4
  jour2/
    etape-05-freelib/     module 5  (go.work + module freelib + provisio)
    etape-06-contexte/ … etape-09-arret/
  jour3/
    etape-10-tests/ … etape-12-perf-secu/
```

Chaque `etape-XX` est l'état complet et compilable du projet à la fin du
module XX. C'est donc à la fois :

- la **correction** du TP du module XX (`TP.md` dedans : énoncé, consignes, vérification),
- la **base** du TP du module XX+1.

`CHANGES.md` dans chaque étape liste ce qui a changé par rapport à la précédente.

Chaque slide qui montre du code indique en bas où il se trouve et comment le
lancer : soit une étape de provisio, soit un programme de `exemples/`.

## Prérequis

- Go 1.27 ou plus récent (`go version`)
- Pour le jour 3 : `mockery` v3 (`go install github.com/vektra/mockery/v3@latest`),
  `dlv` (`go install github.com/go-delve/delve/cmd/dlv@latest`),
  `benchstat` et `govulncheck` (`go install golang.org/x/perf/cmd/benchstat@latest`,
  `go install golang.org/x/vuln/cmd/govulncheck@latest`)

## Lancer

```sh
cd correction/jour1/etape-02-erreurs/provisio && go run ./cmd/provisio
# à partir de l'étape 5, go.work est dans le dossier parent :
cd correction/jour2/etape-05-freelib/provisio && go run ./cmd/provisio
```

```sh
curl -s -X POST localhost:8080/orders \
  -d '{"type":"activation","customer_id":"C42","target":"12 rue des Lilas"}'
curl -s -X POST localhost:8080/orders \
  -d '{"type":"upgrade","customer_id":"C43","target":"CTR-000042"}'
```

`LEXIQUE.md` définit les termes employés dans les slides. Le README de chaque dossier `correction/jourN` détaille les étapes, les variables d'environnement
et les démonstrations.
