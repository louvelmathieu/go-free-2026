# provisio, fil rouge du Jour 2

Suite directe du Jour 1 (`etape-04-fx`). Chaque dossier est l'état complet
et compilable du projet à la fin d'un module.
Chaque étape est aussi la correction du TP de son module : son `TP.md` donne
l'énoncé, le dossier de départ (l'étape précédente), les consignes et la vérification.

| Étape | Module du deck | Ce qui s'ajoute |
|---|---|---|
| `etape-05-freelib` | Modules privés | `freelib` devient un module à part (`logging`, `httpx`), `go.work`, `require freelib v0.1.0`, README de publication (`GOPRIVATE`). À partir d'ici, `go.work` est dans le dossier parent de `provisio` |
| `etape-06-contexte` | context.Context | POST asynchrone (202), timeout global, `DELETE` qui annule, causes, request_id dans le contexte, 3 exemples |
| `etape-07-fuites` | Goroutines qui fuient | système legacy muet, port admin + pprof, compteur de goroutines, test goleak |
| `etape-08-workers` | Patterns de concurrence | Pool de workers, queue bornée, load shedding 503, sémaphores inventaire réseau/système legacy, `sync.Once`, Order concurrency-safe |
| `etape-09-arret` | Arrêt propre | Drainage, délai d'arrêt, `AfterFunc`, `fx.StopTimeout`, `fx.Supply` |

## Lancer

```bash
cd etape-09-arret/provisio
APP_LEGACY_FIX=1 go run ./cmd/provisio
# port public :8080, port admin 127.0.0.1:6060
```

Variables utiles : `APP_ORDER_TIMEOUT=1s`, `APP_LEGACY_FIX=1`,
`APP_LEGACY_HANG_RATE=0`, `APP_STOP_TIMEOUT=2s`, `APP_ENV=prod`.

## Scénarios de démo

```bash
# fuite (étape 7) : sans APP_LEGACY_FIX, 30 commandes upgrade
for i in $(seq 30); do curl -s -o /dev/null -X POST localhost:8080/orders \
  -d '{"type":"upgrade","customer_id":"C1","target":"CTR-000042"}'; done
curl localhost:6060/healthz
curl 'localhost:6060/debug/pprof/goroutine?debug=1' | grep -A3 hang

# race (étapes 6-7) : go run -race ./cmd/provisio puis GET en boucle sur une commande en cours

# load shedding (étape 8) : 60 commandes d'un coup
for i in $(seq 60); do curl -s -o /dev/null -w '%{http_code} ' -X POST localhost:8080/orders \
  -d '{"type":"activation","customer_id":"C1","target":"x"}'; done

# arrêt propre (étape 9) : 20 commandes puis Ctrl-C, avec APP_STOP_TIMEOUT=20s puis 2s
```

Même remarque qu'au Jour 1 : le `replace` vers `../freelib` dans
`provisio/go.mod` n'existe que parce que `gitlab.example.com` est fictif.

## Version de Go

Compilé et testé avec Go 1.27, voir le README du Jour 1.
