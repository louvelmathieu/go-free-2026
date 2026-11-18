// 01-bench/parse_test.go
// CORRECTION : un benchmark par fonction.
//
//	go test ./01-bench -bench . -benchmem
package bench

import "testing"

// sink reçoit le résultat. Sans lui, le compilateur voit que le résultat
// n'est jamais utilisé et peut supprimer l'appel : on mesurerait du vide.
var sink CreateOrderRequest

// Un benchmark est une fonction BenchmarkXxx(b *testing.B).
// b.Loop() (depuis Go 1.24) répète le corps jusqu'à avoir une mesure stable.
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

func BenchmarkParseCSVNoAlloc(b *testing.B) {
	for b.Loop() {
		sink, _ = ParseCSVNoAlloc(CSVIn)
	}
}

// Une optimisation doit donner le même résultat : on le vérifie.
func TestParsersEquivalents(t *testing.T) {
	a, _ := ParseJSON(JSONIn)
	b, _ := ParseCSV(CSVIn)
	c, _ := ParseCSVNoAlloc(CSVIn)
	if a != b || b != c {
		t.Fatalf("résultats différents : %v %v %v", a, b, c)
	}
}
