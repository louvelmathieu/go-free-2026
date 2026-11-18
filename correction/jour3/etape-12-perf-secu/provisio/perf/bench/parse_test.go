package bench

import "testing"

var (
	jsonIn = []byte(`{"type":"activation","customer_id":"C42","target":"12 rue des Lilas, Lyon"}`)
	csvIn  = "activation;C42;12 rue des Lilas, Lyon"
	sink   CreateOrderRequest // empêche le compilateur d'éliminer l'appel
)

// go test ./perf/bench -bench . -benchmem
// b.Loop (Go 1.24) remplace for i := 0; i < b.N; i++ et empêche
// certaines optimisations qui faussaient les mesures.
func BenchmarkParseJSON(b *testing.B) {
	for b.Loop() {
		sink, _ = ParseJSON(jsonIn)
	}
}

func BenchmarkParseCSV(b *testing.B) {
	for b.Loop() {
		sink, _ = ParseCSV(csvIn)
	}
}

func BenchmarkParseCSVNoAlloc(b *testing.B) {
	for b.Loop() {
		sink, _ = ParseCSVNoAlloc(csvIn)
	}
}

func TestParsersEquivalents(t *testing.T) {
	a, _ := ParseJSON(jsonIn)
	b, _ := ParseCSV(csvIn)
	c, _ := ParseCSVNoAlloc(csvIn)
	if a != b || b != c {
		t.Fatalf("résultats différents : %v %v %v", a, b, c)
	}
}
