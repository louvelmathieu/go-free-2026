# Aide-mémoire : sécurité d'un service Go

## Dans le service, ou dans l'infrastructure ?

Sur un vrai projet, un service est rarement exposé directement : une gateway,
un reverse proxy ou un ingress se trouve devant lui.

| Protection | Souvent déjà faite devant le service | À faire dans le service si |
|---|---|---|
| Taille maximale du corps | oui (gateway, nginx, ingress) | rien devant, ou limite propre à une route |
| Délais sur les connexions entrantes | oui | rien devant |
| CORS, en-têtes de sécurité | souvent | l'API est appelée par un navigateur sans gateway |
| TLS, limitation de débit | oui | rien devant |

Ces réglages se justifient dans le service pour un projet personnel, un petit
outil interne, ou un service joint directement sans passer par la gateway.
Vérifier ce que l'infrastructure fait déjà avant de les recopier partout.

## Ce qui reste toujours à la charge du service

- **Les panics des goroutines qu'il lance.** `net/http` rattrape celles des
  handlers, pas les autres : une seule arrête tout le processus. Un `recover`
  en `defer` au début de chaque goroutine lancée à la main.
- **Ses dépendances.** `govulncheck ./...` ne signale que les failles dont le
  code est réellement appelé. À lancer en CI.
- **Les délais sur ses appels sortants.** `http.Client` n'a aucun délai par
  défaut. Aucune gateway ne le fera à sa place.
- **La validation de ce qu'il reçoit**, et des erreurs qui ne révèlent pas
  son fonctionnement interne.
- **Ses bornes de concurrence** (pool, file bornée, sémaphores) : un service
  sans borne tombe avec une simple boucle de requêtes.

## Dépendances

```bash
go install golang.org/x/vuln/cmd/govulncheck@latest
govulncheck ./...
go mod verify               # les modules téléchargés correspondent à go.sum
go list -m -u all           # les mises à jour disponibles
```

`go.sum` se commite et ne se modifie jamais à la main.

## Un serveur HTTP sans rien devant

```go
srv := &http.Server{
    Addr:              ":8080",
    Handler:           mux,
    ReadHeaderTimeout: 5 * time.Second,
    ReadTimeout:       15 * time.Second,
    WriteTimeout:      30 * time.Second,
    IdleTimeout:       60 * time.Second,
}

r.Body = http.MaxBytesReader(w, r.Body, 1<<20)   // dans le handler : 1 Mo maximum
```
