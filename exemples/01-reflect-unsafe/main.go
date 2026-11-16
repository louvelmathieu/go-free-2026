// reflect et unsafe : ce que le compilateur ne vérifie plus.
//
//	go run ./01-reflect-unsafe
package main

import (
	"fmt"
	"reflect"
	"unsafe"
)

type Order struct {
	ID     string
	Status string
}

func (o *Order) Complete() { o.Status = "done" }

func main() {
	o := &Order{ID: "ord_42", Status: "running"}

	// --- reflect : le nom du champ est une string, le compilateur ne la lit pas ---
	func() {
		defer func() {
			if r := recover(); r != nil {
				fmt.Println("reflect, faute de frappe : ça a compilé, et ça panique :", r)
			}
		}()
		v := reflect.ValueOf(o).Elem()
		v.FieldByName("Satus").SetString("done") // "Satus" : panic à l'exécution
	}()

	o.Complete() // la version que le compilateur vérifie
	fmt.Println("via la méthode :", o.Status)

	// --- unsafe : les deux usages sûrs, qui ne font que mesurer ---
	fmt.Println("unsafe.Sizeof(Order{}) =", unsafe.Sizeof(Order{}), "octets")
	fmt.Println("unsafe.Offsetof(Status) =", unsafe.Offsetof(o.Status))

	// --- unsafe : l'usage dangereux, volontairement NON exécuté ---
	// s := "immuable"
	// b := unsafe.Slice(unsafe.StringData(s), len(s))
	// b[0] = 'X' // écrire dans une string : comportement indéfini, souvent un crash
	fmt.Println("(l'écriture dans une string via unsafe est laissée en commentaire : comportement indéfini)")
}
