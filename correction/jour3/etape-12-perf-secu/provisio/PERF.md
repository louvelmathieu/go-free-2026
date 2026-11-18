# Profiler provisio

Le port admin (127.0.0.1:6060) expose pprof. Générer de la charge, puis :

```bash
# CPU, 15 secondes d'échantillonnage, interface web
go tool pprof -http=:8081 http://localhost:6060/debug/pprof/profile?seconds=15

# mémoire : objets vivants (inuse) ou alloués au total (alloc)
go tool pprof -http=:8081 http://localhost:6060/debug/pprof/heap
go tool pprof -sample_index=alloc_space -http=:8081 http://localhost:6060/debug/pprof/heap

# goroutines (jour 2) et blocages
curl 'localhost:6060/debug/pprof/goroutine?debug=1'
go tool pprof http://localhost:6060/debug/pprof/block     # nécessite runtime.SetBlockProfileRate
go tool pprof http://localhost:6060/debug/pprof/mutex     # nécessite runtime.SetMutexProfileFraction

# traceur d'exécution : 5 s, puis l'UI
curl -o trace.out 'localhost:6060/debug/pprof/trace?seconds=5'
go tool trace trace.out
```

Dans l'UI pprof : Top (où le temps passe), Graph, Flame Graph. Source pour la ligne exacte.

## Micro-benchmarks

```bash
go test ./perf/bench -bench . -benchmem
go test ./perf/bench -bench . -benchmem -count 6 > old.txt   # avant
# ... modification ...
go test ./perf/bench -bench . -benchmem -count 6 > new.txt   # après
go install golang.org/x/perf/cmd/benchstat@latest
benchstat old.txt new.txt                                    # différence statistiquement significative ?
```

Règles : `b.Loop()` (depuis Go 1.24) plutôt que `b.N`, un `sink` pour empêcher l'élimination
de code mort, `-count` ≥ 6 et benchstat avant de conclure, jamais sur un portable en batterie.

## Empreinte mémoire des structs

```bash
go run ./perf/alignement
go vet -vettool=$(which fieldalignment) ./...    # golang.org/x/tools/go/analysis/passes/fieldalignment
```

## Taille des binaires

```bash
go build -ldflags="-s -w" ./cmd/provisio    # sans table des symboles ni DWARF : -30 % environ
go version -m provisio                       # ce qui est dedans
```
