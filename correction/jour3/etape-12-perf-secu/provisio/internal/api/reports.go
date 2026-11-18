package api

import (
	"errors"
	"fmt"
	"io/fs"
	"net/http"
	"os"
	"path/filepath"
)

// Reports sert les rapports d'activation stockés dans un dossier.
// Deux implémentations : la naïve, vulnérable au path traversal, et la
// version durcie. Le flag Hardened choisit.
type Reports struct {
	Dir      string
	Hardened bool
}

func (rp *Reports) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	name := r.PathValue("name")

	if !rp.Hardened {
		// VULNÉRABLE : le nom vient du client, on le colle au dossier.
		// GET /reports/..%2F..%2Fetc%2Fpasswd  → /var/provisio/reports/../../etc/passwd
		path := filepath.Join(rp.Dir, name)
		b, err := os.ReadFile(path)
		if err != nil {
			// VERBEUX : le message d'erreur système révèle l'arborescence.
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		w.Write(b) // pas de Content-Type : le navigateur devine (sniffing)
		return
	}

	// DURCI
	// 1. os.Root (Go 1.24) : toute résolution de chemin reste sous Dir,
	//    liens symboliques compris. Impossible d'en sortir.
	root, err := os.OpenRoot(rp.Dir)
	if err != nil {
		http.Error(w, "rapports indisponibles", http.StatusInternalServerError)
		return
	}
	defer root.Close()

	b, err := fs.ReadFile(root.FS(), name)
	if err != nil {
		// 2. Message neutre pour le client, détail dans les logs seulement.
		if errors.Is(err, fs.ErrNotExist) || errors.Is(err, fs.ErrInvalid) {
			http.Error(w, "rapport introuvable", http.StatusNotFound)
			return
		}
		http.Error(w, "rapports indisponibles", http.StatusInternalServerError)
		return
	}

	// 3. Content-Type explicite + nosniff : pas d'interprétation par le navigateur.
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.Header().Set("Content-Length", fmt.Sprint(len(b)))
	w.Write(b)
}
