# Correction du TP 12

Les programmes de ce dossier sont les versions corrigées. Les versions de départ sont dans `tp/12-perf-secu`.
Dans le code corrigé, chaque modification porte le commentaire `CORRECTION`.

## Ce qu'il y a à écrire, en bref

Les numéros de ligne sont ceux des fichiers de `tp/12-perf-secu`.

### 1. `01-bench` : un fichier à créer, une fonction à ajouter

Créer `01-bench/parse_test.go` avec un benchmark par fonction :

```go
package bench

import "testing"

var sink CreateOrderRequest      // empêche le compilateur de supprimer l'appel

func BenchmarkParseJSON(b *testing.B) {
    for b.Loop() {
        sink, _ = ParseJSON(JSONIn)
    }
}

func BenchmarkParseCSV(b *testing.B) {
    for b.Loop() {
        sink, _ = ParseCSV(CSVIn)
    }
}
```

Puis ajouter dans `parse.go` la version sans allocation, et son benchmark :

```go
func ParseCSVNoAlloc(s string) (CreateOrderRequest, error) {
    i := strings.IndexByte(s, ';')
    if i < 0 {
        return CreateOrderRequest{}, errors.New("3 champs attendus")
    }
    j := strings.IndexByte(s[i+1:], ';')
    if j < 0 {
        return CreateOrderRequest{}, errors.New("3 champs attendus")
    }
    j += i + 1
    return CreateOrderRequest{Type: s[:i], CustomerID: s[i+1 : j], Target: s[j+1:]}, nil
}
```

### 2. `02-alignement/main.go` : une struct à ajouter

```go
type OrderCompact struct {
    CreatedAt int64
    Priority  int64
    Weight    float64
    Retries   int32
    Done      bool
    Cancelled bool
}
```

Les mêmes six champs, du plus gros au plus petit : 32 octets au lieu de 48.

### 3. `03-govulncheck` : une commande, aucune ligne de code

```sh
go get golang.org/x/text@latest
```

`go.mod`, ligne 5, passe de `golang.org/x/text v0.3.5` à une version corrigée. Le fichier Go ne change pas.

### 4. `04-serveur/main.go` : trois corrections

**A. Ligne 36, le corps illimité : une ligne à ajouter avant, et le cas d'erreur à traiter**

```go
r.Body = http.MaxBytesReader(w, r.Body, 1<<20)    // ajoutée : 1 Mo maximum

data, err := io.ReadAll(r.Body)
if err != nil {
    var tropGros *http.MaxBytesError               // ajoutées : répondre 413
    if errors.As(err, &tropGros) {
        http.Error(w, "fichier trop volumineux (1 Mo maximum)", http.StatusRequestEntityTooLarge)
        return
    }
    http.Error(w, "lecture impossible", http.StatusBadRequest)
    return
}
```

**B. Ligne 51, aucun délai : une ligne à remplacer par un serveur configuré**

```go
// Avant
log.Fatal(http.ListenAndServe(":8080", mux))

// Après
srv := &http.Server{
    Addr:              ":8080",
    Handler:           mux,
    ReadHeaderTimeout: 5 * time.Second,
    ReadTimeout:       15 * time.Second,
    WriteTimeout:      30 * time.Second,
    IdleTimeout:       60 * time.Second,
}
log.Fatal(srv.ListenAndServe())
```

**C. Lignes 46 et 56, la panic en arrière-plan : une ligne à remplacer, cinq à ajouter**

```go
// Ligne 46, avant
go envoyerNotification(nil)
// Après
go envoyerNotification(map[string]int{})

// Au début de envoyerNotification, à ajouter
defer func() {
    if r := recover(); r != nil {
        log.Println("panic dans envoyerNotification :", r)
    }
}()
```

Pour voir toutes les différences d'un coup : `diff -ru tp/12-perf-secu correction/jour3/etape-12-perf-secu`.

---

Le détail de chaque exercice :

## 1. `01-bench` : lire un benchmark

```
$ go test ./01-bench -bench . -benchmem
BenchmarkParseJSON-2           1995867     646.5 ns/op     48 B/op     1 allocs/op
BenchmarkParseCSV-2           22269234      54.67 ns/op    48 B/op     1 allocs/op
BenchmarkParseCSVNoAlloc-2    76006689      15.96 ns/op     0 B/op     0 allocs/op
```

Les chiffres changent d'une machine à l'autre, les rapports restent proches.

- Le grand nombre après le nom est le nombre de répétitions choisi par Go. Ce n'est pas un score.
- `ns/op`, `B/op`, `allocs/op` : le temps, les octets et le nombre d'allocations par appel. Les trois se lisent dans le même sens, plus c'est bas, mieux c'est.
- JSON est 12 fois plus lent que le CSV, pour la même mémoire.
- Supprimer la seule allocation du CSV (le slice créé par `strings.Split`) divise le temps par trois.

**À faire remarquer** : 0,6 microseconde par commande ne se verra jamais sur un service qui en traite cent par seconde. Le benchmark sert aussi à décider de ne pas optimiser.

## 2. `02-alignement` : l'ordre des champs

```
$ go run ./02-alignement
somme des champs : 1 + 8 + 4 + 8 + 1 + 8 = 30 octets
OrderNaif        : 48 octets
OrderCompact     : 32 octets
```

Le compilateur place chaque champ sur un multiple de sa taille. Un `bool` suivi d'un `int64` laisse 7 octets vides. `Offsetof` montre les trous : `Done` est à 0, `CreatedAt` à 8.

**À faire remarquer** : ça compte avec des millions d'instances en mémoire, pas pour une struct de configuration. Le linter `fieldalignment` le signale.

## 3. `03-govulncheck` : une dépendance vulnérable

```
$ govulncheck ./03-govulncheck
Vulnerability #1: GO-2022-1059
    Denial of service via crafted Accept-Language header in
    golang.org/x/text/language
  Module: golang.org/x/text
    Found in: golang.org/x/text@v0.3.5
    Fixed in: golang.org/x/text@v0.3.8
    Example traces found:
      #1: 03-govulncheck/main.go:24:46: 03.LangueDuClient calls language.ParseAcceptLanguage

Your code is affected by 2 vulnerabilities from 1 module.
```

Comment le lire :

- **Le titre** : un en-tête `Accept-Language` fabriqué exprès peut bloquer le programme. Or cet en-tête vient du client.
- **Found in / Fixed in** : la version utilisée, et la première version corrigée.
- **Example traces** : la ligne de notre code qui mène à la fonction vulnérable. C'est ce qui distingue govulncheck d'un simple inventaire : il ne signale que ce qu'on appelle réellement.

Après `go get golang.org/x/text@latest` :

```
$ govulncheck ./03-govulncheck
No vulnerabilities found.
```

**À faire remarquer** : le programme marchait parfaitement. Aucun test n'aurait trouvé ce problème. Et la dernière ligne du rapport mentionne aussi des failles présentes dans le module mais jamais appelées : celles-là ne nous concernent pas.

## 4. `04-serveur` : trois faiblesses

### A. Le corps illimité

```
$ head -c 50000000 /dev/zero | curl -X POST --data-binary @- localhost:8080/import
50000000 octets importés
```

Le serveur a chargé 50 Mo en mémoire (119 Mo occupés d'après `ps`). Avec 2 Go, il tombe. Après correction, il répond 413 et sa mémoire reste à 9 Mo.

### B. Le client lent

```
$ go run ./04-serveur/client-lent
après 30s, le serveur garde toujours la connexion ouverte
```

`http.ListenAndServe` n'a aucun délai. Après correction :

```
le serveur a répondu au bout de 15s : "HTTP/1.1 400 Bad Request"
```

C'est le `ReadTimeout` de 15 secondes.

### C. La panic en arrière-plan

```
$ curl localhost:8080/notifier
notification en cours
$ curl localhost:8080/sante
curl: (7) Failed to connect
```

Le serveur est mort. `net/http` rattrape les panics des handlers, mais pas celles des goroutines qu'on lance soi-même : une seule suffit à arrêter tout le processus. C'est la leçon du TP 11, exercice 1, vue côté serveur.

## Ce qu'il faut dire sur la portée de l'exercice 4

Sur un vrai projet, le service est rarement exposé directement. Une gateway, un reverse proxy ou un ingress se trouve devant, et pose déjà une taille maximale de corps et des délais. Refaire ces réglages dans chaque service n'est alors pas forcément utile.

Ces protections se justifient dans le service quand il n'y a rien devant : un projet personnel, un petit outil interne, un service appelé directement par d'autres services sans passer par la gateway.

Ce qui reste toujours à la charge du service, quelle que soit l'infrastructure :

- rattraper les panics des goroutines qu'il lance,
- surveiller ses dépendances (govulncheck),
- poser un délai sur les appels qu'il fait lui-même vers l'extérieur.
