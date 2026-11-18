# Aide-mémoire : mesurer la performance

## Benchmarks

```bash
go test ./01-bench -bench . -benchmem
```

| Colonne | Ce que c'est | On la veut |
|---|---|---|
| le nombre après le nom | combien de fois Go a répété le code | ce n'est pas un score |
| `ns/op` | le temps moyen d'un appel, en nanosecondes | basse |
| `B/op` | la mémoire allouée par appel, en octets | basse |
| `allocs/op` | le nombre d'allocations par appel | basse, idéalement 0 |

Comparer avant et après une modification :

```bash
go test ./01-bench -bench . -benchmem -count 6 > avant.txt
# ... modification ...
go test ./01-bench -bench . -benchmem -count 6 > apres.txt
go install golang.org/x/perf/cmd/benchstat@latest
benchstat avant.txt apres.txt        # l'écart est-il réel, ou du bruit ?
```

Règles : `b.Loop()` (depuis Go 1.24), une variable `sink` pour que le
compilateur ne supprime pas l'appel, plusieurs séries et benchstat avant de
conclure, et une machine calme.

## Profiler un benchmark

```bash
go test ./01-bench -bench . -cpuprofile cpu.out      # 1. fabrique le fichier cpu.out
go tool pprof -top cpu.out | head -12                # 2. le lit et affiche le classement
go tool pprof -http=:8081 cpu.out                    # ou dans le navigateur : Top, Flame Graph
```

Rien n'écoute et rien ne tourne en parallèle : `cpu.out` est un fichier.

## Profiler un service qui tourne

Il faut trois terminaux. Avec provisio (`correction/jour3/etape-10-tests/provisio`),
dont le port d'administration 6060 expose pprof :

```bash
# Terminal 1 : le service
go run ./cmd/provisio

# Terminal 2 : pprof va chercher un profil sur le port 6060 (il attend 15 s),
# puis ouvre le port 8081 pour l'afficher dans le navigateur
go tool pprof -http=:8081 'http://localhost:6060/debug/pprof/profile?seconds=15'

# Terminal 3, pendant ces 15 secondes : envoyer des requêtes, sinon le profil est vide
```

Les autres profils, sur la même adresse :

```bash
go tool pprof -http=:8081 http://localhost:6060/debug/pprof/heap      # mémoire
curl 'localhost:6060/debug/pprof/goroutine?debug=1'                   # goroutines
curl -o trace.out 'localhost:6060/debug/pprof/trace?seconds=5' && go tool trace trace.out
```
