# Débugger provisio avec Delve

## Installer

```bash
go install github.com/go-delve/delve/cmd/dlv@latest
```

GoLand et VS Code (extension Go) l'utilisent en dessous. `.vscode/launch.json`
contient les configurations prêtes.

## Les quatre modes

| Commande | Quand |
|---|---|
| `dlv debug ./cmd/provisio` | compiler et lancer sous debugger |
| `dlv test ./internal/orchestrator -- -test.run TestOrchestratorSuite/TestActivation_Succes` | un test précis |
| `dlv attach <pid>` | un process déjà lancé (serveur figé, deadlock de mutex) |
| `dlv exec ./provisio` | un binaire déjà compilé (avec `-gcflags=all="-N -l"` pour garder les variables lisibles) |

Sur un serveur distant : `dlv --headless --listen=:2345 --api-version=2 attach <pid>`
puis la configuration « distant » de VS Code, ou `dlv connect host:2345`.

## Les commandes qu'on utilise vraiment

```
break orchestrator.go:120          # breakpoint (fichier:ligne, ou pkg.Func)
break (*Orchestrator).step
condition 1 name == "legacy"          # ne s'arrêter que si la condition est vraie
continue  / next  / step  / stepout
locals  / args  / print order.status  / print *order
goroutines                         # toutes les goroutines et où elles sont
goroutine 42 bt                    # la stack de la goroutine 42
goroutine 42                       # s'y placer, puis locals, print...
watch -w s.available               # s'arrêter quand la valeur change
trace fake.(*Legacy).Provision        # tracer les appels sans s'arrêter
```

## Le bug du jour

```bash
go test ./debug/03-fuite-stock                 # FAIL : des équipements ont fui
dlv test ./debug/03-fuite-stock
(dlv) break (*Stock).Release
(dlv) continue                               # compter les passages...
(dlv) break stock.go:60                      # la branche billing
(dlv) condition 2 attempt == 3
(dlv) continue
(dlv) print s.available                      # avant et après : l'équipement n'est pas rendue
```

## Sans debugger

- `kill -QUIT <pid>` : le runtime écrit toutes les stacks sur stderr. Le
  réflexe pour un service figé en prod sans Delve.
- `GOTRACEBACK=all` : une panic affiche toutes les goroutines, pas
  seulement la fautive.
- `curl localhost:6060/debug/pprof/goroutine?debug=2` : les mêmes stacks,
  via le port admin, sans toucher au process.

## Lire une stack

1. Le message de panic ou l'état `[chan receive]`, `[semacquire]`, `[select]`.
2. La première ligne qui est dans **votre** code (pas `runtime/`, pas une lib).
3. `created by` : qui a lancé cette goroutine.
4. Pour un deadlock de mutex : chercher deux goroutines en `[sync.Mutex.Lock]`
   sur le même objet, ou une seule qui le prend deux fois.
