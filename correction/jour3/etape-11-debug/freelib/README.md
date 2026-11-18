# freelib

Module commun de l'équipe : logging, middleware HTTP, client HTTP configuré.

## Publier une version

```bash
git tag v0.1.0
git push origin v0.1.0
```

Un module Go, c'est un dépôt git et un tag semver. Rien de plus à "publier".
Pour un v2 et au-delà, le chemin du module doit finir par `/v2`
(`gitlab.example.com/free/freelib/v2`).

## Consommer depuis un dépôt privé

Go passe par proxy.golang.org et sum.golang.org par défaut. Ils ne voient pas
vos dépôts privés, il faut leur dire de ne pas essayer :

```bash
go env -w GOPRIVATE=gitlab.example.com/free/*
```

`GOPRIVATE` implique `GONOPROXY` et `GONOSUMDB` pour ces chemins : Go clone
directement en git. Il faut donc que git sache s'authentifier sans prompt :

```bash
# HTTPS avec token
git config --global url."https://oauth2:${GITLAB_TOKEN}@gitlab.example.com/".insteadOf "https://gitlab.example.com/"

# ou SSH
git config --global url."git@gitlab.example.com:".insteadOf "https://gitlab.example.com/"
```

Puis, dans le service :

```bash
go get gitlab.example.com/free/freelib@v0.1.0
```

## Développer provisio contre une version locale de freelib

Deux mécanismes :

### `go.work` (recommandé)

À la racine du dossier qui contient les deux clones :

```
go 1.27

use (
    ./freelib
    ./provisio
)
```

Go résout freelib depuis le dossier local. `go.work` ne se commite pas dans
provisio (il est propre au poste), et il est ignoré par `go mod tidy`, par
`go build` en CI si `GOWORK=off`, et par `go install` d'un module distant.

### `replace` dans go.mod

```
replace gitlab.example.com/free/freelib => ../freelib
```

Fonctionne aussi, mais le `replace` est dans go.mod, donc facile à commiter
par erreur : la CI casse (le chemin n'existe pas) ou pire, elle passe avec
une version de freelib qui n'est pas celle taguée. À réserver aux forks
durables, pas au développement du quotidien.

## Vérifier qu'on ne dépend pas du local

```bash
GOWORK=off go build ./...
```

Si ça échoue, c'est que le code dépend d'un changement de freelib pas encore
tagué : taguer freelib, puis `go get` la nouvelle version dans provisio.
