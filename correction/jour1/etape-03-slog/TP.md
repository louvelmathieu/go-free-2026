# TP 3 : Des logs qu'on peut suivre

## Objectif

Remplacer le package log par slog, et faire en sorte qu'on puisse suivre une commande de bout en bout avec un seul identifiant.

## Point de départ

Ouvrez `correction/jour1/etape-02-erreurs` (la correction du TP 2). Copiez-le dans un dossier de travail à vous, ou travaillez directement dedans sur une branche.

```sh
cp -r correction/jour1/etape-02-erreurs mon-tp-03
cd mon-tp-03
```

## Consignes

1. Construire un *slog.Logger dans main : handler texte en dev, JSON quand APP_ENV=prod, avec les attributs service et env.
2. Écrire un middleware HTTP qui génère un request_id, le met dans un logger dérivé, et logue méthode, chemin, statut et durée à la fin.
3. Faire voyager ce logger dans le contexte (clé de type non exporté, fonctions WithLogger et From) et l'utiliser dans l'orchestrateur.
4. Implémenter slog.LogValuer sur Order pour loguer id, type et statut sans exposer la struct.

## Vérification

- APP_ENV=prod go run ./cmd/provisio puis un POST : chaque ligne est un JSON avec le même request_id.
- Le log de fin d'étape montre le statut courant de la commande, pas celui du début (piège du With).

## Un conseil

Un logger par requête, jamais un logger global dans le cœur. Et un attribut qui change souvent (statut, compteur) se passe à chaque appel, pas dans With.

## Correction

Ce dossier (`correction/jour1/etape-03-slog`) est la correction : l'état complet du projet à la fin du module. `CHANGES.md` liste ce qui a changé par rapport à l'étape précédente, et un `diff -r` entre les deux dossiers montre chaque ligne.

```sh
diff -ru correction/jour1/etape-02-erreurs correction/jour1/etape-03-slog | less
```

Dépôt : https://github.com/louvelmathieu/go-free-2026
