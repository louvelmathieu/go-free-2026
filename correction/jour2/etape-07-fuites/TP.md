# TP 7 : Trouver et boucher une fuite de goroutines

## Objectif

Reproduire une fuite causée par un système externe qui ne répond pas, la voir avec pprof, la corriger, et écrire le test qui empêche son retour.

## Point de départ

Ouvrez `correction/jour2/etape-06-contexte` (la correction du TP 6). Copiez-le dans un dossier de travail à vous, ou travaillez directement dedans sur une branche.

```sh
cp -r correction/jour2/etape-06-contexte mon-tp-07
cd mon-tp-07   # go.work à la racine, le code dans provisio/
```

## Consignes

1. Donner au fake Legacy un HangRate : une proportion d'appels qui ne répondent jamais, en ignorant le contexte (c'est le bug).
2. Exposer net/http/pprof sur un serveur d'administration séparé, en écoute sur 127.0.0.1:6060 seulement.
3. Envoyer une trentaine de commandes upgrade, lire /debug/pprof/goroutine?debug=1 et retrouver la stack qui s'accumule.
4. Corriger : le fake respecte ctx.Done() (option RespectContext, activée par APP_LEGACY_FIX=1 pour la démo), et l'étape sort en timeout au lieu de pendre.
5. Écrire legacy_test.go avec goleak.VerifyNone : il doit échouer sur la version muette et passer sur la version corrigée.

## Vérification

- Avant correction : le nombre de goroutines monte à chaque salve et ne redescend pas.
- Après : il revient à sa valeur de départ, et go test -race ./internal/clients/... passe.

## Un conseil

Une goroutine qui attend sur un channel ou un appel réseau sans surveiller ctx.Done() est une fuite en attente. Le select à deux branches est le réflexe à prendre.

## Correction

Ce dossier (`correction/jour2/etape-07-fuites`) est la correction : l'état complet du projet à la fin du module. `CHANGES.md` liste ce qui a changé par rapport à l'étape précédente, et un `diff -r` entre les deux dossiers montre chaque ligne.

```sh
diff -ru correction/jour2/etape-06-contexte correction/jour2/etape-07-fuites | less
```

Dépôt : https://github.com/louvelmathieu/go-free-2026
