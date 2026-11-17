# TP 6 : Rendre la commande asynchrone et annulable

## Objectif

Le POST répond tout de suite, la commande s'exécute en arrière-plan avec un délai global, et le client peut l'annuler.

## Point de départ

Ouvrez `correction/jour2/etape-05-freelib` (la correction du TP 5). Copiez-le dans un dossier de travail à vous, ou travaillez directement dedans sur une branche.

```sh
cp -r correction/jour2/etape-05-freelib mon-tp-06
cd mon-tp-06   # go.work à la racine, le code dans provisio/
```

## Consignes

1. POST /orders enregistre la commande, répond 202, et lance l'orchestration dans une goroutine avec un contexte dérivé du serveur (pas de la requête).
2. Poser un délai global par commande avec context.WithTimeoutCause(ctx, cfg.OrderTimeout, domain.ErrOrderTimeout).
3. Ajouter DELETE /orders/{id} : le store garde la fonction d'annulation (WithCancelCause) et l'appelle avec ErrCancelledByClient.
4. Dans l'orchestrateur : vérifier context.Cause(ctx) entre les étapes, et faire la compensation (Release) avec context.WithoutCancel.
5. writeError traduit ErrOrderTimeout et ErrCancelledByClient, et GET montre le statut cancelled.

## Vérification

- POST renvoie 202, GET passe de running à done en une seconde ou deux.
- POST puis DELETE immédiat : GET montre cancelled, et le log montre la cause exacte.
- APP_ORDER_TIMEOUT=300ms : les commandes finissent en failed avec « délai global de la commande dépassé ».

## Un conseil

Le contexte de la requête meurt avec la réponse. Tout ce qui doit survivre au 202 se dérive d'un autre parent, et c'est ce parent qu'on annulera à l'arrêt.

## Correction

Ce dossier (`correction/jour2/etape-06-contexte`) est la correction : l'état complet du projet à la fin du module. `CHANGES.md` liste ce qui a changé par rapport à l'étape précédente, et un `diff -r` entre les deux dossiers montre chaque ligne.

```sh
diff -ru correction/jour2/etape-05-freelib correction/jour2/etape-06-contexte | less
```

Dépôt : https://github.com/louvelmathieu/go-free-2026
