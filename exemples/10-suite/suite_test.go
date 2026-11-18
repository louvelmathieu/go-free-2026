// 10-suite/suite_test.go
// Une suite testify qui affiche l'ordre dans lequel ses méthodes sont appelées.
//
//	go test ./10-suite -v
package suite_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/suite"
)

// DemoSuite : une struct qui embarque suite.Suite.
type DemoSuite struct {
	suite.Suite
	compteur int // un état partagé, remis à zéro avant chaque test
}

// Les quatre méthodes ci-dessous ne sont pas des mots-clés de Go. Ce sont des
// noms que testify cherche sur la struct : si la méthode existe, il l'appelle.
// Le nom doit être exact, majuscules comprises.

// SetupSuite : UNE fois, avant le premier test.
func (s *DemoSuite) SetupSuite() { fmt.Println("SetupSuite     (une fois, avant tout)") }

// SetupTest : avant CHAQUE test.
func (s *DemoSuite) SetupTest() {
	s.compteur = 0
	fmt.Println("  SetupTest    (avant chaque test) : compteur remis à 0")
}

// TearDownTest : après CHAQUE test.
func (s *DemoSuite) TearDownTest() { fmt.Println("  TearDownTest (après chaque test)") }

// TearDownSuite : UNE fois, après le dernier test.
func (s *DemoSuite) TearDownSuite() { fmt.Println("TearDownSuite  (une fois, après tout)") }

// Les tests : toute méthode dont le nom commence par Test.
func (s *DemoSuite) TestA() {
	s.compteur++
	fmt.Println("    TestA : compteur =", s.compteur)
	s.Equal(1, s.compteur)
}

func (s *DemoSuite) TestB() {
	s.compteur++
	fmt.Println("    TestB : compteur =", s.compteur)
	s.Equal(1, s.compteur, "TestB ne voit pas ce que TestA a fait : SetupTest est repassé")
}

// Le seul point d'entrée que go test connaît. Sans cette fonction, rien ne
// s'exécute, et rien ne le signale.
func TestDemoSuite(t *testing.T) {
	suite.Run(t, new(DemoSuite))
}
