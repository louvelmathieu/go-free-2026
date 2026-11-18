// 04-serveur/client-lent/main.go
// Un client volontairement lent : il annonce un corps de 1000 octets, en
// envoie un seul, puis ne dit plus rien. Il mesure combien de temps le
// serveur garde la connexion ouverte.
//
//	go run ./04-serveur/client-lent
package main

import (
	"fmt"
	"net"
	"time"
)

func main() {
	conn, err := net.Dial("tcp", "127.0.0.1:8080")
	if err != nil {
		fmt.Println("le serveur ne répond pas :", err)
		return
	}
	defer conn.Close()

	debut := time.Now()
	fmt.Fprint(conn, "POST /import HTTP/1.1\r\nHost: localhost\r\nContent-Length: 1000\r\n\r\nx")
	fmt.Println("en-têtes et 1 octet envoyés, puis plus rien. On attend 30 secondes au plus...")

	conn.SetReadDeadline(time.Now().Add(30 * time.Second))
	buf := make([]byte, 64)
	n, err := conn.Read(buf)
	duree := time.Since(debut).Round(time.Second)

	if err != nil && n == 0 {
		if ne, ok := err.(net.Error); ok && ne.Timeout() {
			fmt.Printf("après %v, le serveur garde toujours la connexion ouverte\n", duree)
			fmt.Println("Avec quelques milliers de clients comme celui-ci, il ne peut plus servir personne.")
			return
		}
		fmt.Printf("le serveur a fermé la connexion au bout de %v\n", duree)
		return
	}
	fmt.Printf("le serveur a répondu au bout de %v : %q\n", duree, string(buf[:min(n, 24)]))
}
