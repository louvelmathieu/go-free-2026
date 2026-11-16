# TP 2 : Des erreurs qu'on peut lire

## Objectif

Faire remonter les erreurs d'étape avec assez de contexte pour que l'API puisse répondre juste, sans que l'orchestrateur connaisse HTTP.

## Point de départ

Ouvrez `correction/jour1/etape-01-structure` (la correction du TP 1). Copiez-le dans un dossier de travail à vous, ou travaillez directement dedans sur une branche.

```sh
cp -r correction/jour1/etape-01-structure mon-tp-02
cd mon-tp-02
```

## Consignes

1. Dans domain : des sentinelles (ErrNotFound, ErrNotEligible, ErrNoStock, ErrTransient), une ValidationError typée et une StepError{OrderID, Step, Err} qui implémente Unwrap.
2. Dans l'orchestrateur : une seule fonction step() qui pose le timeout, wrappe l'erreur dans StepError et marque le timeout comme ErrTransient.
3. Sur l'upgrade : si le legacy échoue après la réservation d'équipement, libérer l'équipement et conserver les deux erreurs avec errors.Join.
4. Dans api : une fonction writeError qui traduit avec errors.Is / errors.As vers 400, 404, 409, 422 ou 502, avec un corps JSON (error, field, step, retryable).

## Vérification

- Une commande vers « hors-zone » renvoie 422 et {"error":"adresse non éligible"} au lieu d'un 500.
- Un JSON sans customer_id renvoie 400 avec le champ en faute.
- go vet ./... ne signale aucun %w mal placé.

## Un conseil

Le test à se poser pour chaque erreur : qui va la lire, et que va-t-il en faire ? Une sentinelle pour décider, un type pour transporter des détails, un simple %w quand on ne fait qu'ajouter du contexte.

## Correction

Ce dossier (`correction/jour1/etape-02-erreurs`) est la correction : l'état complet du projet à la fin du module. `CHANGES.md` liste ce qui a changé par rapport à l'étape précédente, et un `diff -r` entre les deux dossiers montre chaque ligne.

```sh
diff -ru correction/jour1/etape-01-structure correction/jour1/etape-02-erreurs | less
```

Dépôt : https://github.com/louvelmathieu/go-free-2026
