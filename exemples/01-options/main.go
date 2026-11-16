// Trois façons de configurer un constructeur. La dernière est le pattern
// "functional options".
//
//	go run ./01-options
package main

import (
	"fmt"
	"time"
)

type Orchestrator struct {
	name        string
	stepTimeout time.Duration
	retries     int
}

//  1. Paramètres positionnels : illisible à l'appel, et chaque nouveau réglage
//     casse tous les appels existants.
func NewPositional(name string, stepTimeout time.Duration, retries int) *Orchestrator {
	return &Orchestrator{name: name, stepTimeout: stepTimeout, retries: retries}
}

//  2. Struct de config : lisible, extensible, mais la valeur zéro est ambiguë.
//     StepTimeout: 0 veut dire "pas renseigné" ou "zéro seconde" ?
type Config struct {
	Name        string
	StepTimeout time.Duration
	Retries     int
}

func NewFromConfig(cfg Config) *Orchestrator {
	if cfg.StepTimeout == 0 { // on devine...
		cfg.StepTimeout = 5 * time.Second
	}
	return &Orchestrator{name: cfg.Name, stepTimeout: cfg.StepTimeout, retries: cfg.Retries}
}

// 3. Functional options : une Option est une fonction qui modifie la struct.
type Option func(*Orchestrator)

func WithStepTimeout(d time.Duration) Option { return func(o *Orchestrator) { o.stepTimeout = d } }
func WithRetries(n int) Option               { return func(o *Orchestrator) { o.retries = n } }

func New(name string, opts ...Option) *Orchestrator {
	o := &Orchestrator{name: name, stepTimeout: 5 * time.Second} // les défauts d'abord
	for _, opt := range opts {
		opt(o) // puis les options, dans l'ordre donné
	}
	return o
}

func main() {
	fmt.Printf("%+v\n", *NewPositional("a", 2*time.Second, 3))                     // 2s, c'est quoi ?
	fmt.Printf("%+v\n", *NewFromConfig(Config{Name: "b"}))                         // timeout deviné
	fmt.Printf("%+v\n", *New("c"))                                                 // défauts
	fmt.Printf("%+v\n", *New("d", WithRetries(3), WithStepTimeout(2*time.Second))) // nommé, ordre libre
}
