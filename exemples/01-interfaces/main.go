// Une interface en Go : implicite, petite, déclarée par celui qui l'utilise.
//
//	go run ./01-interfaces
package main

import (
	"context"
	"fmt"
	"strings"
)

// --- côté consommateur : l'orchestrateur dit ce dont il a besoin ---

// Inventory est déclarée ici, dans le package qui l'appelle. Une méthode suffit.
type Inventory interface {
	CheckEligibility(ctx context.Context, address string) (bool, error)
}

func activate(ctx context.Context, inv Inventory, address string) error {
	ok, err := inv.CheckEligibility(ctx, address)
	if err != nil {
		return err
	}
	if !ok {
		return fmt.Errorf("%s : non éligible", address)
	}
	fmt.Println(address, ": éligible, on continue")
	return nil
}

// --- côté fournisseur : un type concret, qui ne mentionne jamais l'interface ---

type FakeInventory struct{}

func (FakeInventory) CheckEligibility(_ context.Context, address string) (bool, error) {
	return !strings.Contains(address, "hors-zone"), nil
}

// Un second type, écrit par quelqu'un d'autre, sans se concerter : il a la
// méthode, donc il satisfait l'interface. C'est ça, "implicite".
type AlwaysYes struct{}

func (AlwaysYes) CheckEligibility(context.Context, string) (bool, error) { return true, nil }

// Vérification à la compilation, si on veut être explicite quelque part :
var _ Inventory = FakeInventory{}
var _ Inventory = AlwaysYes{}

func main() {
	ctx := context.Background()
	// La fonction accepte une interface : n'importe quel type qui a la méthode passe.
	activate(ctx, FakeInventory{}, "12 rue des Lilas")
	fmt.Println(activate(ctx, FakeInventory{}, "hameau hors-zone"))
	activate(ctx, AlwaysYes{}, "hameau hors-zone")
}
