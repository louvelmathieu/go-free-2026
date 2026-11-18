# TP 12 : Mesurer, puis durcir

## Objectif

Quatre petits exercices indépendants de provisio : lire un benchmark, mesurer une struct, trouver une dépendance vulnérable, et mettre en difficulté un serveur HTTP avant de le corriger.

## Point de départ

Ouvrez `tp/12-perf-secu` (les programmes à corriger, indépendants de provisio). Copiez-le dans un dossier de travail à vous, ou travaillez directement dedans sur une branche.

```sh
cp -r tp/12-perf-secu mon-tp-12
cd mon-tp-12
```

## Consignes

1. 01-bench : écrire un benchmark par fonction (b.Loop), lancer avec -benchmem, lire ns/op, B/op et allocs/op. Puis écrire ParseCSVNoAlloc, sans allocation, et mesurer le gain.
2. 02-alignement : lancer, repérer les octets perdus avec les positions affichées, écrire OrderCompact avec les mêmes champs dans un autre ordre.
3. 03-govulncheck : lancer govulncheck ./03-govulncheck, lire le rapport (quelle faille, quelle version la corrige, quelle ligne y mène), corriger et relancer.
4. 04-serveur : lancer le serveur, faire les trois essais décrits en tête du fichier (corps énorme, client lent, panic en arrière-plan), puis corriger.

## Vérification

- ParseCSVNoAlloc affiche 0 allocs/op, et govulncheck répond « No vulnerabilities found ».
- Le serveur répond 413 au corps de 50 Mo, ferme le client lent au bout de 15 s, et répond toujours sur /sante après /notifier.

## Un conseil

Pas d'optimisation sans mesure, pas de correction sans avoir vu le problème. Dans les quatre exercices, on observe d'abord, on modifie ensuite.

## Correction

Ce dossier (`correction/jour3/etape-12-perf-secu`) est la correction : l'état complet du projet à la fin du module. `CHANGES.md` liste ce qui a changé par rapport à l'étape précédente, et un `diff -r` entre les deux dossiers montre chaque ligne.

```sh
diff -ru tp/12-perf-secu correction/jour3/etape-12-perf-secu | less
```

Dépôt : https://github.com/louvelmathieu/go-free-2026
