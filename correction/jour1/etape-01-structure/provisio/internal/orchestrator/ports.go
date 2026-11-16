package orchestrator

import "context"

// Les interfaces sont déclarées là où elles sont consommées, et elles sont
// petites : une implémentation fake ou un mock tient en quelques lignes.

// Inventory pilote l'inventaire réseau : éligibilité et ressources.
type Inventory interface {
	CheckEligibility(ctx context.Context, address string) (bool, error)
	ReserveResource(ctx context.Context, address string) (port string, err error)
}

// Configurator configure l'équipement du client.
type Configurator interface {
	Configure(ctx context.Context, customerID, port string) error
}

// Stock gère le stock d'équipements.
type Stock interface {
	Reserve(ctx context.Context) (ref string, err error)
	Release(ctx context.Context, ref string) error
}

// LegacyClient provisionne le service dans le système de provisioning historique.
type LegacyClient interface {
	Provision(ctx context.Context, target, ref string) error
}

// Billing ouvre la facturation d'un service.
type Billing interface {
	Open(ctx context.Context, customerID, line string) error
}
