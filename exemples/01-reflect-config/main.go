// 01-reflect-config/main.go
// reflect, un cas concret : remplir N'IMPORTE QUELLE struct de config
// depuis les variables d'environnement, en lisant ses tags.
//
//	go run ./01-reflect-config
package main

import (
	"fmt"
	"os"
	"reflect"
	"strconv"
)

// Deux structs différentes, que LoadEnv ne connaît pas à l'avance.
type HTTPConfig struct {
	Addr    string `env:"APP_ADDR"`
	Workers int    `env:"APP_WORKERS"`
}

type LegacyConfig struct {
	URL     string `env:"LEGACY_URL"`
	Retries int    `env:"LEGACY_RETRIES"`
}

// --- Sans reflect : une fonction par struct, une ligne par champ ---
// C'est le bon choix tant qu'on a UNE config dans UN service.
func loadHTTPALaMain() HTTPConfig {
	workers, _ := strconv.Atoi(os.Getenv("APP_WORKERS"))
	return HTTPConfig{Addr: os.Getenv("APP_ADDR"), Workers: workers}
}

// --- Avec reflect : une seule fonction pour toutes les structs ---
// dst est un pointeur vers une struct quelconque.
func LoadEnv(dst any) error {
	v := reflect.ValueOf(dst) // la valeur, vue à l'exécution
	if v.Kind() != reflect.Pointer || v.Elem().Kind() != reflect.Struct {
		return fmt.Errorf("LoadEnv: attendu un pointeur de struct, reçu %T", dst)
	}
	v = v.Elem()  // la struct pointée (modifiable, car on a le pointeur)
	t := v.Type() // son type : c'est lui qui connaît les champs et les tags

	for i := 0; i < t.NumField(); i++ {
		name := t.Field(i).Tag.Get("env") // le tag `env:"..."` du champ i
		raw, ok := os.LookupEnv(name)
		if name == "" || !ok {
			continue // pas de tag ou variable absente : on laisse le zéro
		}
		field := v.Field(i)
		switch field.Kind() { // on ne sait écrire que les types prévus ici
		case reflect.String:
			field.SetString(raw)
		case reflect.Int:
			n, err := strconv.Atoi(raw)
			if err != nil {
				return fmt.Errorf("LoadEnv: %s=%q n'est pas un entier", name, raw)
			}
			field.SetInt(int64(n))
		default:
			return fmt.Errorf("LoadEnv: type %s non géré (champ %s)", field.Kind(), t.Field(i).Name)
		}
	}
	return nil
}

func main() {
	os.Setenv("APP_ADDR", ":8080")
	os.Setenv("APP_WORKERS", "8")
	os.Setenv("LEGACY_URL", "http://legacy.local")
	os.Setenv("LEGACY_RETRIES", "3")

	fmt.Printf("à la main      : %+v\n", loadHTTPALaMain())

	var h HTTPConfig
	var l LegacyConfig
	fmt.Println("LoadEnv(&h)    :", LoadEnv(&h), fmt.Sprintf("%+v", h))
	fmt.Println("LoadEnv(&l)    :", LoadEnv(&l), fmt.Sprintf("%+v", l))

	// Le prix : ce que le compilateur ne vérifie plus.
	fmt.Println("LoadEnv(h)     :", LoadEnv(h)) // oubli du & : ça compile, erreur à l'exécution
	os.Setenv("APP_WORKERS", "huit")
	fmt.Println("valeur fausse  :", LoadEnv(&h))
	fmt.Println()
	fmt.Println("C'est exactement ce que font encoding/json, fx et sheriff :")
	fmt.Println("du code écrit une fois, qui marche sur des types qu'il n'a jamais vus.")
}
