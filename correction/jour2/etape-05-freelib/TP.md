# TP 5 : Sortir le code commun dans freelib

## Objectif

Extraire les deux packages réutilisables de provisio dans un module à part, comme une librairie d'équipe, et continuer à développer les deux côte à côte.

## Point de départ

Ouvrez `correction/jour1/etape-04-fx` (la correction du TP 4). Copiez-le dans un dossier de travail à vous, ou travaillez directement dedans sur une branche.

```sh
cp -r correction/jour1/etape-04-fx mon-tp-05
cd mon-tp-05   # go.work à la racine, le code dans provisio/
```

## Consignes

1. Créer un module freelib (gitlab.example.com/free/freelib) à côté de provisio, et y déplacer tels quels internal/logging et internal/httpx.
2. Mettre à jour les imports dans internal/app, internal/api et internal/orchestrator : provisio ne garde plus que du code métier.
3. Ajouter un go.work à la racine qui use ./freelib ./provisio, et un require freelib v0.1.0 dans provisio/go.mod.
4. Écrire dans freelib/README.md la procédure de publication (tag v0.1.0, GOPRIVATE, GONOSUMDB) et quand préférer go.work à replace.

## Vérification

- Depuis la racine : go build ./... compile les deux modules, et go run ./cmd/provisio se comporte comme hier.
- go list -m all dans provisio montre freelib, et une modification locale de freelib est prise en compte sans tag.
- Le go.work n'est pas commité dans provisio (il décrit votre poste, pas le projet).

## Un conseil

Le module est l'unité de version, pas le package. Un seul module freelib avec plusieurs packages suffit tant que tout est publié ensemble.

## Correction

Ce dossier (`correction/jour2/etape-05-freelib`) est la correction : l'état complet du projet à la fin du module. `CHANGES.md` liste ce qui a changé par rapport à l'étape précédente, et un `diff -r` entre les deux dossiers montre chaque ligne.

```sh
diff -ru correction/jour1/etape-04-fx correction/jour2/etape-05-freelib | less
```

Dépôt : https://github.com/louvelmathieu/go-free-2026
