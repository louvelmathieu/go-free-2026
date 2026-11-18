// internal/clients/spy/spy.go
// Package spy fournit des doubles de test écrits à la main pour les cinq
// systèmes externes. Un spy fait deux choses, et rien d'autre :
//
//  1. il répond ce qu'on lui a dit de répondre (ses champs publics),
//  2. il note chaque appel reçu, avec ses arguments, pour qu'un test puisse
//     vérifier après coup « as-tu bien été appelé, et avec quoi ? ».
//
// Aucune bibliothèque, aucune génération de code : tout est lisible ici.
package spy

import (
	"context"
	"sync"
)

// Call est un appel reçu : le nom de la méthode et ses arguments (sans le ctx).
type Call struct {
	Method string
	Args   []any
}

// Recorder est le carnet de notes commun à tous les spies.
// Il est protégé par un mutex : les workers appellent en parallèle.
type Recorder struct {
	mu    sync.Mutex
	calls []Call
}

// record note un appel. Chaque méthode d'un spy commence par là.
func (r *Recorder) record(method string, args ...any) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.calls = append(r.calls, Call{Method: method, Args: args})
}

// Calls renvoie une copie de tous les appels reçus, dans l'ordre.
func (r *Recorder) Calls() []Call {
	r.mu.Lock()
	defer r.mu.Unlock()
	return append([]Call(nil), r.calls...)
}

// Count dit combien de fois une méthode a été appelée.
func (r *Recorder) Count(method string) int {
	n := 0
	for _, c := range r.Calls() {
		if c.Method == method {
			n++
		}
	}
	return n
}

// Args renvoie les arguments du premier appel à une méthode, ou nil.
func (r *Recorder) Args(method string) []any {
	for _, c := range r.Calls() {
		if c.Method == method {
			return c.Args
		}
	}
	return nil
}

// Inventory remplace l'inventaire réseau.
type Inventory struct {
	Recorder
	Eligible bool   // ce que CheckEligibility répond
	Port     string // ce que ReserveResource répond
	Err      error  // si non nil, les deux méthodes échouent avec
}

// NewInventory : par défaut l'adresse est éligible et un port est disponible.
func NewInventory() *Inventory { return &Inventory{Eligible: true, Port: "RES-01/2"} }

func (s *Inventory) CheckEligibility(_ context.Context, address string) (bool, error) {
	s.record("CheckEligibility", address)
	return s.Eligible, s.Err
}

func (s *Inventory) ReserveResource(_ context.Context, address string) (string, error) {
	s.record("ReserveResource", address)
	return s.Port, s.Err
}

// Configurator remplace le configurateur d'équipement.
type Configurator struct {
	Recorder
	Err error
}

func (s *Configurator) Configure(_ context.Context, customerID, port string) error {
	s.record("Configure", customerID, port)
	return s.Err
}

// Stock remplace le stock d'équipements.
type Stock struct {
	Recorder
	Ref        string // la référence que Reserve renvoie
	ReserveErr error
	ReleaseErr error
}

// NewStock : par défaut une référence est toujours disponible.
func NewStock() *Stock { return &Stock{Ref: "EQ-0001"} }

func (s *Stock) Reserve(_ context.Context) (string, error) {
	s.record("Reserve")
	return s.Ref, s.ReserveErr
}

func (s *Stock) Release(_ context.Context, ref string) error {
	s.record("Release", ref)
	return s.ReleaseErr
}

// Legacy remplace le système de provisioning historique.
type Legacy struct {
	Recorder
	Err error
	// Do, s'il est renseigné, remplace la réponse : le spy exécute ce code.
	// Sert à simuler un système lent qui attend le contexte.
	Do func(ctx context.Context) error
}

func (s *Legacy) Provision(ctx context.Context, target, ref string) error {
	s.record("Provision", target, ref)
	if s.Do != nil {
		return s.Do(ctx)
	}
	return s.Err
}

// Billing remplace la facturation.
type Billing struct {
	Recorder
	Err error
}

func (s *Billing) Open(_ context.Context, customerID, line string) error {
	s.record("Open", customerID, line)
	return s.Err
}
