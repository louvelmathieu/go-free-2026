# TP 10 : Tester l'orchestrateur sans ses systèmes

## Objectif

Couvrir l'orchestrateur, le handler et l'application fx entière avec des spies écrits à la main, sans bibliothèque de mock.

## Point de départ

Ouvrez `correction/jour2/etape-09-arret` (la correction du TP 9). Copiez-le dans un dossier de travail à vous, ou travaillez directement dedans sur une branche.

```sh
cp -r correction/jour2/etape-09-arret mon-tp-10
cd mon-tp-10   # go.work à la racine, le code dans provisio/
```

## Consignes

1. Écrire un spy par système externe dans internal/clients/spy : des champs pour régler la réponse, un carnet des appels reçus (méthode et arguments).
2. Tester l'orchestrateur avec ces spies : succès, non éligible, échec legacy avec compensation, timeout d'étape, annulation. Vérifier les arguments transmis et ce qui ne doit pas être appelé.
3. Tester le handler avec httptest, un spy du pool écrit dans le test et le vrai store en mémoire : 202, 400, 404, 503 file pleine.
4. Remonter toute l'application avec fxtest.New, remplacer les cinq systèmes avec fx.Replace, récupérer le handler avec fx.Populate, envoyer un POST et lire les carnets des spies.

## Vérification

- make test et make race verts, make cover au-dessus de 70 % sur orchestrator et api.
- En retirant le fx.As d'un fx.Replace, le test de bout en bout échoue : le fake est resté en place.

## Un conseil

Un spy par interface consommée, pas par struct. Si une interface est pénible à doubler, c'est elle qu'il faut réduire, pas le test.

## Correction

Ce dossier (`correction/jour3/etape-10-tests`) est la correction : l'état complet du projet à la fin du module. `CHANGES.md` liste ce qui a changé par rapport à l'étape précédente, et un `diff -r` entre les deux dossiers montre chaque ligne.

```sh
diff -ru correction/jour2/etape-09-arret correction/jour3/etape-10-tests | less
```

Dépôt : https://github.com/louvelmathieu/go-free-2026
