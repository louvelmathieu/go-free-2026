# TP 9 : S'arrêter sans perdre de commande

## Objectif

Sur SIGTERM, finir ce qui est en cours dans un délai donné, annuler le reste proprement, et refuser ce qui arrive entre-temps.

## Point de départ

Ouvrez `correction/jour2/etape-08-workers` (la correction du TP 8). Copiez-le dans un dossier de travail à vous, ou travaillez directement dedans sur une branche.

```sh
cp -r correction/jour2/etape-08-workers mon-tp-09
cd mon-tp-09   # go.work à la racine, le code dans provisio/
```

## Consignes

1. Pool.Stop : fermer la file aux nouvelles soumissions (ErrStopping → 503), attendre les workers jusqu'à StopTimeout, puis annuler le contexte des commandes restantes avec une cause explicite.
2. Lire la config avant fx.New pour passer fx.StopTimeout(cfg.StopTimeout) et fx.Supply(cfg).
3. Garder une marge (500 ms) entre le délai du pool et celui de fx, pour que l'arrêt forcé soit loggé avant que fx abandonne.
4. Utiliser context.AfterFunc pour tracer chaque commande annulée par l'arrêt.

## Vérification

- 20 commandes puis SIGTERM avec APP_STOP_TIMEOUT=20s : toutes terminées.
- La même chose avec APP_STOP_TIMEOUT=2s : quelques terminées, le reste cancelled avec la cause « arrêt du service avant la fin de la commande », et le processus sort quand même.

## Un conseil

Un arrêt propre, c'est trois décisions : qu'est-ce qu'on refuse, qu'est-ce qu'on finit, et au bout de combien de temps on lâche. Si une des trois est implicite, elle sera prise par un kill -9.

## Correction

Ce dossier (`correction/jour2/etape-09-arret`) est la correction : l'état complet du projet à la fin du module. `CHANGES.md` liste ce qui a changé par rapport à l'étape précédente, et un `diff -r` entre les deux dossiers montre chaque ligne.

```sh
diff -ru correction/jour2/etape-08-workers correction/jour2/etape-09-arret | less
```

Dépôt : https://github.com/louvelmathieu/go-free-2026
