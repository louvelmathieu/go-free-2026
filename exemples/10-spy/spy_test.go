// 10-spy/spy_test.go
// Un spy écrit à la main, en entier sur une page : il répond ce qu'on lui a
// réglé et il note ce qu'il a reçu. Aucune bibliothèque de mock.
//
//	go test ./10-spy -v
package spy

import (
	"context"
	"errors"
	"sync"
	"testing"
)

// --- Le code à tester ---

// LegacyClient est le système externe dont dépend Activer.
type LegacyClient interface {
	Provision(ctx context.Context, target, ref string) error
}

// Stock réserve et rend des équipements.
type Stock interface {
	Reserve(ctx context.Context) (string, error)
	Release(ctx context.Context, ref string) error
}

// Activer réserve un équipement puis provisionne. Si le provisioning
// échoue, l'équipement doit être rendu : c'est ce qu'on veut vérifier.
func Activer(ctx context.Context, st Stock, lg LegacyClient, target string) error {
	ref, err := st.Reserve(ctx)
	if err != nil {
		return err
	}
	if err := lg.Provision(ctx, target, ref); err != nil {
		_ = st.Release(ctx, ref) // la compensation
		return err
	}
	return nil
}

// --- Les spies ---

// Call est un appel reçu : le nom de la méthode et ses arguments.
type Call struct { // (1)
	Method string
	Args   []any
}

// SpyLegacy remplace le système legacy dans un test.
type SpyLegacy struct {
	mu    sync.Mutex
	calls []Call // (2) le carnet de notes
	Err   error  // (3) la réponse, réglable depuis le test
}

// Provision a la même signature que l'interface : le spy prend la place du vrai. (4)
func (s *SpyLegacy) Provision(_ context.Context, target, ref string) error {
	s.mu.Lock()
	s.calls = append(s.calls, Call{"Provision", []any{target, ref}}) // (5) on note
	s.mu.Unlock()
	return s.Err // (6) on répond ce qui a été réglé
}

// SpyStock : le même principe, pour deux méthodes.
type SpyStock struct {
	mu    sync.Mutex
	calls []Call
	Ref   string
}

func (s *SpyStock) Reserve(context.Context) (string, error) {
	s.mu.Lock()
	s.calls = append(s.calls, Call{"Reserve", nil})
	s.mu.Unlock()
	return s.Ref, nil
}

func (s *SpyStock) Release(_ context.Context, ref string) error {
	s.mu.Lock()
	s.calls = append(s.calls, Call{"Release", []any{ref}})
	s.mu.Unlock()
	return nil
}

// --- Les tests : régler, lancer, vérifier ---

func TestActiver_Succes(t *testing.T) {
	stock := &SpyStock{Ref: "EQ-0001"}
	legacy := &SpyLegacy{} // Err vaut nil : le legacy répond « tout va bien »

	err := Activer(context.Background(), stock, legacy, "CTR-42")

	if err != nil {
		t.Fatalf("erreur inattendue : %v", err)
	}
	// Le legacy a été appelé une fois, avec la cible et la référence réservée.
	if len(legacy.calls) != 1 {
		t.Fatalf("Provision appelé %d fois, attendu 1", len(legacy.calls))
	}
	if got := legacy.calls[0].Args; got[0] != "CTR-42" || got[1] != "EQ-0001" {
		t.Errorf("Provision reçu %v, attendu [CTR-42 EQ-0001]", got)
	}
	// Rien n'a été rendu au stock : un seul appel, Reserve.
	if len(stock.calls) != 1 {
		t.Errorf("le stock a reçu %v, attendu seulement Reserve", stock.calls)
	}
}

func TestActiver_LegacyEnPanne_RendLEquipement(t *testing.T) {
	stock := &SpyStock{Ref: "EQ-0001"}
	legacy := &SpyLegacy{Err: errors.New("système legacy indisponible")} // on règle la panne

	err := Activer(context.Background(), stock, legacy, "CTR-42")

	if err == nil {
		t.Fatal("une erreur était attendue")
	}
	// La compensation : le dernier appel au stock doit être Release("EQ-0001").
	last := stock.calls[len(stock.calls)-1]
	if last.Method != "Release" || last.Args[0] != "EQ-0001" {
		t.Errorf("dernier appel au stock : %v, attendu Release(EQ-0001)", last)
	}
	t.Logf("appels reçus par le stock : %v", stock.calls)
}
