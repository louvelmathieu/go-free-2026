# Étape 8 : patterns de concurrence

- `internal/worker` : pool de taille fixe (`APP` `Workers: 8`) devant une
  queue bornée (`QueueDepth: 32`). `Submit` est non bloquant : file pleine
  = `ErrQueueFull` = 503 + `Retry-After` (load shedding). `sync.Once` rend
  `Start` idempotent.
- `internal/clients/limited` : wrappers à sémaphore
  (`golang.org/x/sync/semaphore`) autour de l'inventaire réseau (3 connexions) et du
  système legacy (20). `Acquire(ctx)` respecte l'annulation.
- `domain.Order` devient concurrency-safe (RWMutex) : la race détectée à
  l'étape 7 avec `-race` disparaît.
- fx : les fakes inventaire réseau/système legacy sont fournis sous leur type concret, les
  wrappers `limited` sous les interfaces (`fx.As`).
- Démo : 60 POST d'un coup → 40 × 202 puis 20 × 503 ; `/healthz` montre
  `pending` qui redescend à 0.
