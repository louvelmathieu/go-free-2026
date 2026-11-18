# Étape 11 : debug

Cette étape ne contient pas provisio. Ce sont trois petits programmes
indépendants, chacun avec un bug à trouver. Les versions avec les bugs sont
dans `tp/11-debug`, ce dossier contient les versions corrigées.

- `01-panic` : une panic dans une goroutine (écriture dans une map nil).
- `02-deadlock` : un blocage que Go détecte (channel) et un qu'il ne détecte
  pas (`-mutex`), à inspecter avec `kill -QUIT` ou `dlv attach`.
- `03-fuite-stock` : un test qui échoue, à comprendre avec `dlv test` et un
  breakpoint conditionnel.
- `CORRECTION.md` : pour chaque bug, ce qu'on voit, comment on le trouve, la
  cause et la correction.
- `DEBUG.md` : aide-mémoire Delve et lecture de stacks.
