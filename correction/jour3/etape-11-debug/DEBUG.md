# Aide-mémoire : débugger avec Delve

## Installer

```bash
go install github.com/go-delve/delve/cmd/dlv@latest
dlv version
```

GoLand l'embarque (bouton Debug). VS Code l'installe avec l'extension Go, et
`tp/11-debug/.vscode/launch.json` contient les configurations prêtes.

## Les quatre façons d'entrer

| Commande | Quand |
|---|---|
| `dlv debug ./01-panic` | compiler et lancer un programme sous debugger |
| `dlv test ./03-fuite-stock` | lancer les tests d'un package sous debugger |
| `dlv attach <pid>` | un programme déjà lancé (figé, deadlock de mutex) |
| `dlv exec ./binaire` | un binaire déjà compilé (avec `-gcflags=all="-N -l"` pour garder les variables lisibles) |

Le `./` est obligatoire. Sur un serveur distant : `dlv --headless --listen=:2345 --api-version=2 attach <pid>`,
puis la configuration « distant » de VS Code, ou `dlv connect host:2345`.

## Les commandes qu'on utilise vraiment

```
break main.main                    # breakpoint sur une fonction
break stock.go:54                  # ou sur fichier:ligne
condition 1 attempt == 3           # ne s'arrêter que si la condition est vraie
continue  / next  / step  / stepout
locals  / args  / print s.available  / print *s
goroutines                         # toutes les goroutines et où elles sont
goroutine 1 bt                     # la stack de la goroutine 1
watch -w s.available               # s'arrêter quand la valeur change
quit                               # ou Ctrl-D. Ctrl-C ne quitte pas
```

L'invite `(dlv)` attend une commande : le programme ne tourne pas tant qu'on n'a pas tapé `continue`.

## Sans debugger

- `kill -QUIT <pid>` : le runtime écrit toutes les stacks sur stderr. Le
  réflexe pour un service figé en prod sans Delve.
- `GOTRACEBACK=all` : une panic affiche toutes les goroutines, pas
  seulement la fautive.
- `curl localhost:6060/debug/pprof/goroutine?debug=2` : les mêmes stacks,
  sur un service qui expose pprof (provisio, port d'administration).

## Lire une stack

1. Le message de panic ou l'état `[chan receive]`, `[semacquire]`, `[select]`.
2. La première ligne qui est dans **votre** code (pas `runtime/`, pas une lib).
3. `created by` : qui a lancé cette goroutine.
4. Pour un deadlock de mutex : chercher deux goroutines en `[sync.Mutex.Lock]`
   sur le même objet, ou une seule qui le prend deux fois.
