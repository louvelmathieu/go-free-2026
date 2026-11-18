# Étape 12 : performance et sécurité

Cette étape ne contient pas provisio. Ce sont quatre petits exercices
indépendants. Les versions de départ sont dans `tp/12-perf-secu`, ce dossier
contient les versions corrigées. La dernière version de provisio est donc
celle de l'étape 10.

- `01-bench` : les benchmarks (`b.Loop`, `-benchmem`) de deux parseurs, et
  une version sans allocation.
- `02-alignement` : l'effet de l'ordre des champs sur la taille d'une struct.
- `03-govulncheck` : une dépendance vulnérable trouvée par govulncheck,
  corrigée par une montée de version.
- `04-serveur` : un serveur HTTP par défaut, corrigé sur trois points (taille
  du corps, délais, panic dans une goroutine).
- `CORRECTION.md` : ce qu'il y a à écrire pour chaque exercice, puis le détail.
- `PERF.md` : aide-mémoire benchmarks et pprof.
- `SECURITE.md` : ce qui relève du service, ce qui relève de l'infrastructure.
