# TP 11 : Trois bugs, un debugger

## Objectif

Diagnostiquer un panic, un deadlock et un bug de stock avec les outils de Go et Delve, sans ajouter un seul Println.

## Point de départ

Ouvrez `tp/11-debug` (les programmes à corriger, indépendants de provisio). Copiez-le dans un dossier de travail à vous, ou travaillez directement dedans sur une branche.

```sh
cp -r tp/11-debug mon-tp-11
cd mon-tp-11
```

## Consignes

1. 01-panic : lancer, lire la trace, retrouver la ligne fautive et la goroutine qui l'a exécutée, corriger.
2. 02-deadlock : lancer le cas 1 et lire le message. Lancer le cas 2 (-mutex), envoyer SIGQUIT, lire les stacks et expliquer qui attend quoi avant de corriger.
3. 03-fuite-stock : go test ./03-fuite-stock échoue. Avec dlv test (ou le bouton Debug de l'IDE), poser un breakpoint conditionnel sur attempt == 3, suivre s.available pas à pas, trouver le chemin qui ne rend pas l'équipement.

## Vérification

- go run ./01-panic et go run ./02-deadlock (avec et sans -mutex) se terminent sans erreur, go test ./03-fuite-stock passe.
- Vous savez attacher dlv au programme 02-deadlock figé et afficher la stack de sa goroutine bloquée.

## Un conseil

Le breakpoint conditionnel remplace 90 % des Println. Et la trace d'un panic se lit de haut en bas : la première ligne de votre code est la bonne.

## Correction

Ce dossier (`correction/jour3/etape-11-debug`) est la correction : l'état complet du projet à la fin du module. `CHANGES.md` liste ce qui a changé par rapport à l'étape précédente, et un `diff -r` entre les deux dossiers montre chaque ligne.

```sh
diff -ru tp/11-debug correction/jour3/etape-11-debug | less
```

Dépôt : https://github.com/louvelmathieu/go-free-2026
