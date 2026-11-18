# TP 8 : Borner la concurrence

## Objectif

Remplacer la goroutine par commande par un pool de workers et une queue bornée, et protéger les systèmes externes avec des sémaphores.

## Point de départ

Ouvrez `correction/jour2/etape-07-fuites` (la correction du TP 7). Copiez-le dans un dossier de travail à vous, ou travaillez directement dedans sur une branche.

```sh
cp -r correction/jour2/etape-07-fuites mon-tp-08
cd mon-tp-08   # go.work à la racine, le code dans provisio/
```

## Consignes

1. Écrire internal/worker : un Pool avec N workers lisant une queue bornée, Start protégé par sync.Once, Submit non bloquant qui renvoie ErrQueueFull.
2. Le handler traduit ErrQueueFull en 503 avec un en-tête Retry-After.
3. Créer internal/clients/limited : des décorateurs Inventory et Legacy qui prennent un jeton dans un x/sync/semaphore.Weighted avant l'appel (3 pour l'inventaire, 20 pour le legacy).
4. Protéger Order avec un sync.RWMutex : l'API lit pendant que le worker écrit, go test -race doit être vert.
5. Rendre Workers, QueueDepth, InventoryMaxConns et LegacyMaxConns configurables.

## Vérification

- 60 POST lancés en parallèle : une partie en 202, le reste en 503 avec Retry-After, aucun 500.
- go test -race ./... passe, et un test du pool vérifie le load shedding.

## Un conseil

Un channel bufferisé n'est pas une file de travail : il n'a ni refus explicite ni arrêt. Le pool rend ces deux décisions visibles.

## Correction

Ce dossier (`correction/jour2/etape-08-workers`) est la correction : l'état complet du projet à la fin du module. `CHANGES.md` liste ce qui a changé par rapport à l'étape précédente, et un `diff -r` entre les deux dossiers montre chaque ligne.

```sh
diff -ru correction/jour2/etape-07-fuites correction/jour2/etape-08-workers | less
```

Dépôt : https://github.com/louvelmathieu/go-free-2026
