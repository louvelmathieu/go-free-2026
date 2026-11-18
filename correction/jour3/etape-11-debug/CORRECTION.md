# Correction du TP 11

Les programmes de ce dossier sont les versions corrigées. Les versions avec les bugs sont dans `tp/11-debug`.
Dans le code corrigé, chaque ligne modifiée porte le commentaire `// CORRECTION`.

## Les lignes à modifier, en bref

Quatre corrections, huit lignes en tout. Les numéros de ligne sont ceux des fichiers de `tp/11-debug`.

### 1. `01-panic/main.go`, ligne 33 : une ligne à remplacer

```go
// Avant
var reg Registry

// Après
reg := Registry{byID: map[string]*Order{}}
```

### 2. `02-deadlock/main.go`, ligne 36 : une ligne à remplacer (cas 1, le channel)

```go
// Avant
results <- "done"

// Après
go func() { results <- "done" }()
```

### 3. `02-deadlock/main.go`, lignes 50 et 55 à 59 : cinq lignes (cas 2, le mutex)

```go
// Avant
func (s *Stock) Reserve() {
    s.mu.Lock()
    defer s.mu.Unlock()
    if s.count == 0 {
        s.Refill()
    }
    s.count--
}

func (s *Stock) Refill() {
    s.mu.Lock()
    defer s.mu.Unlock()
    s.count = 10
}

// Après
func (s *Stock) Reserve() {
    s.mu.Lock()
    defer s.mu.Unlock()
    if s.count == 0 {
        s.refillLocked()          // modifiée : n'essaie plus de reprendre le lock
    }
    s.count--
}

func (s *Stock) Refill() {
    s.mu.Lock()
    defer s.mu.Unlock()
    s.refillLocked()              // modifiée
}

func (s *Stock) refillLocked() {  // ajoutée : suppose le lock déjà pris
    s.count = 10
}
```

### 4. `03-fuite-stock/stock.go`, ligne 55 : une ligne à ajouter

```go
// Avant
if !billingOK(attempt) {
    return fmt.Errorf("activation %d: %w", attempt, errBilling)
}

// Après
if !billingOK(attempt) {
    s.Release(ref)                // ajoutée
    return fmt.Errorf("activation %d: %w", attempt, errBilling)
}
```

Pour voir toutes les différences d'un coup : `diff -ru tp/11-debug correction/jour3/etape-11-debug`.

---

Le détail de chaque bug, dans l'ordre du TP :

## 1. `01-panic` : écrire dans une map nil

**Ce qu'on voit**

```
$ go run ./01-panic
démarrage
panic: assignment to entry in nil map

goroutine 7 [running]:
main.(*Registry).Add(...)
        .../01-panic/main.go:23
main.process(...)
        .../01-panic/main.go:29
created by main.main in goroutine 1
        .../01-panic/main.go:35
```

**Comment on le trouve**

- Ligne 1 : le message dit tout, on écrit dans une map qui vaut nil.
- La première ligne de notre code est `main.go:23`, dans `Registry.Add` : `r.byID[o.ID] = o`.
- `created by main.main` : la goroutine a été lancée par le `go process(...)` de `main`.

**La cause** : `var reg Registry` crée une struct dont le champ `byID` vaut nil. On peut lire une map nil, mais y écrire fait planter.

**La correction** : créer la map au moment de déclarer le registre, `reg := Registry{byID: map[string]*Order{}}`. Dans un vrai projet, on rangerait cette ligne dans un constructeur `NewRegistry()` pour que personne ne puisse l'oublier.

**À faire remarquer** : la panic a lieu dans une goroutine secondaire, et pourtant tout le programme s'arrête. Une panic non rattrapée dans n'importe quelle goroutine tue le processus entier.

## 2. `02-deadlock` : deux blocages

### Cas 1, détecté par Go

```
$ go run ./02-deadlock
fatal error: all goroutines are asleep - deadlock!

goroutine 1 [chan send]:
main.channelDeadlock()
        .../02-deadlock/main.go:36
```

- `[chan send]` : la goroutine 1 est bloquée sur un envoi.
- Ligne 36 : `results <- "done"`. Le channel n'est pas bufferisé, donc l'envoi attend un lecteur. Le lecteur est la ligne d'après, dans la même goroutine : elle ne sera jamais atteinte.
- Go le signale parce que toutes les goroutines sont bloquées.

**La correction** : faire l'envoi dans une autre goroutine (ou donner une capacité de 1 au channel).

### Cas 2, pas détecté

```
$ go run ./02-deadlock -mutex
tick : le programme est vivant
tick : le programme est vivant
...
```

Le programme ne plante pas et n'avance pas. Go ne dit rien, parce qu'une goroutine (le ticker) tourne encore.

**Comment on le trouve**, dans un second terminal :

```
$ kill -QUIT $(pgrep -f 02-deadlock)
```

Go affiche toutes les stacks. Celle qui compte :

```
goroutine 1 [sync.Mutex.Lock]:
...
sync.(*Mutex).Lock(...)
main.(*Stock).Refill(...)
main.(*Stock).Reserve(...)
main.mutexDeadlock()
```

- `[sync.Mutex.Lock]` : la goroutine 1 attend un lock.
- Elle l'attend dans `Refill`, appelée par `Reserve`. Or `Reserve` a déjà pris ce lock. La goroutine s'attend elle-même.

Avec Delve : `dlv attach $(pgrep -f 02-deadlock)`, puis `goroutines` et `goroutine 1 bt` donnent la même information.

**La cause** : un `sync.Mutex` n'est pas réentrant. La même goroutine ne peut pas le prendre deux fois.

**La correction** : une méthode interne `refillLocked()` qui ne prend pas le lock, appelée par `Reserve`. La méthode publique `Refill` prend le lock puis appelle `refillLocked`.

## 3. `03-fuite-stock` : un équipement jamais rendu

**Ce qu'on voit**

```
$ go test ./03-fuite-stock
--- FAIL: TestActivation_StockRevientAZero
    stock = 2, attendu 5 (5 activations réussies) : des équipements ont fui
```

Il manque 3 équipements. Le test fait échouer la facturation aux essais 3, 6 et 9 : trois échecs, trois équipements perdus. C'est l'indice.

**Comment on le trouve**

```
$ dlv test ./03-fuite-stock
(dlv) break Activate
(dlv) condition 1 attempt == 3
(dlv) continue              # s'arrête à l'entrée d'Activate, au 3e essai
(dlv) next
(dlv) next                  # on est sur s.Reserve(ref)
(dlv) print s.available     # 8 : deux équipements déjà partis aux essais 1 et 2
(dlv) next
(dlv) print s.available     # 7 : celui-ci est réservé
(dlv) next                  # provision réussit
(dlv) next                  # billingOK renvoie false, on arrive sur le return
(dlv) print s.available     # toujours 7 : l'équipement n'a pas été rendu
```

Le breakpoint conditionnel évite de s'arrêter aux essais 1 et 2, qui se passent bien.

**La cause** : dans `Activate`, la branche « provisioning en échec » appelle `s.Release(ref)`, mais la branche « facturation en échec » fait son `return` sans rendre l'équipement.

**La correction** : ajouter `s.Release(ref)` avant ce `return`.

**Le lien avec le reste de la formation** : c'est le même comportement que le test `TestUpgrade_EchecLegacy_LibereLEquipement` du module 10 protège dans provisio. Le matin on a écrit le test qui empêche ce bug, l'après-midi on apprend à le trouver quand le test n'existe pas.

## Vérifier que tout est corrigé

```
$ go run ./01-panic
$ go run ./02-deadlock
$ go run ./02-deadlock -mutex
$ go test ./03-fuite-stock
```

Les quatre commandes se terminent sans erreur.
