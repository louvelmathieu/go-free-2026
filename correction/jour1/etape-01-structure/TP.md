# TP 1 : Découper le monolithe

## Objectif

Partir d'un service écrit dans un seul fichier et le réorganiser en packages qui ne se marchent pas dessus, sans changer son comportement.

## Point de départ

Ouvrez `tp/00-monolithe` (le service en un seul fichier). Copiez-le dans un dossier de travail à vous, ou travaillez directement dedans sur une branche.

```sh
cp -r tp/00-monolithe mon-tp-01
cd mon-tp-01
```

## Consignes

1. Créer les packages internal/domain, internal/orchestrator, internal/api, internal/store et internal/clients/fake, et y déplacer le code de main.go.
2. Rendre les champs de Order privés : on ne construit un Order que par NewOrder (qui valide), on ne le modifie que par ses méthodes.
3. Déclarer les interfaces des systèmes externes dans le package orchestrator, celui qui les consomme, et les garder à une ou deux méthodes.
4. Donner à l'orchestrateur un constructeur à functional options (WithInventory, WithStock, WithStepTimeout...).
5. Le store doit être utilisable à partir de sa valeur vierge : var s store.Memory suffit.

## Vérification

- go build ./... et go vet ./... passent.
- curl -X POST localhost:8080/orders -d '{"type":"activation","customer_id":"C42","target":"12 rue des Lilas"}' renvoie la commande en status done.
- main.go ne contient plus que du câblage : pas de logique métier.

## Un conseil

Ne cherchez pas le découpage parfait. Un package par responsabilité, et des imports qui vont tous dans le même sens : api → orchestrator → domain. Si un import remonte, c'est qu'un type est au mauvais endroit.

## Correction

Ce dossier (`correction/jour1/etape-01-structure`) est la correction : l'état complet du projet à la fin du module. `CHANGES.md` liste ce qui a changé par rapport à l'étape précédente, et un `diff -r` entre les deux dossiers montre chaque ligne.

```sh
diff -ru tp/00-monolithe correction/jour1/etape-01-structure | less
```

Dépôt : https://github.com/louvelmathieu/go-free-2026
