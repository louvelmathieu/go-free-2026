# TP 4 : Laisser fx assembler le service

## Objectif

Remplacer le câblage manuel de main par un graphe fx, et savoir lire ce que fx dit quand ça ne démarre pas.

## Point de départ

Ouvrez `correction/jour1/etape-03-slog` (la correction du TP 3). Copiez-le dans un dossier de travail à vous, ou travaillez directement dedans sur une branche.

```sh
cp -r correction/jour1/etape-03-slog mon-tp-04
cd mon-tp-04
```

## Consignes

1. Créer internal/config (une struct Config lue depuis l'environnement) et internal/app avec un fx.Module par responsabilité : Config, Logging, Clients, Orchestrator, HTTP.
2. Fournir les fakes sous leurs interfaces avec fx.Annotate + fx.As, et grouper les dépendances de l'orchestrateur dans une struct fx.In.
3. Démarrer et arrêter le http.Server avec fx.Lifecycle (OnStart non bloquant, OnStop qui appelle Shutdown).
4. Déplacer le middleware HTTP et NewID de internal/api vers un package internal/httpx, et donner à logging.New le nom du service en paramètre : on prépare l'extraction de demain.
5. Brancher fxevent.SlogLogger pour voir le graphe dans vos logs, et réduire main.go à fx.New(app.Modules...).Run().
6. Puis provoquer les trois pièges du dossier pieges/ : Provide sans Invoke, type fourni deux fois, cycle. Lire chaque erreur avant de la corriger.

## Vérification

- go run ./cmd/provisio démarre, répond aux POST, et Ctrl-C arrête proprement avec les hooks dans l'ordre inverse.
- Retirer un fx.Provide : l'erreur de démarrage nomme le type manquant et la fonction qui le demandait.

## Un conseil

Provide déclare comment fabriquer, Invoke déclenche. Rien de ce qui est Provide ne s'exécute tant qu'un Invoke, directement ou non, ne le demande.

## Correction

Ce dossier (`correction/jour1/etape-04-fx`) est la correction : l'état complet du projet à la fin du module. `CHANGES.md` liste ce qui a changé par rapport à l'étape précédente, et un `diff -r` entre les deux dossiers montre chaque ligne.

```sh
diff -ru correction/jour1/etape-03-slog correction/jour1/etape-04-fx | less
```

Dépôt : https://github.com/louvelmathieu/go-free-2026
