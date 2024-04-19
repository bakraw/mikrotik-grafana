package main

import (
	"fmt"
	"log"
	"time"

	probing "github.com/prometheus-community/pro-bing"
)

// Ping une adresse IP pour vérifier son état.
// Prend en entrée un adresse IP (string) et renvoie le statut (int, up = 1 et down = 0).
// On peut changer le nombre de paquets à envoyer et la durée avant time out.
func probeIP(IPaddr string) int {
	var up int

	// Configuration du ping
	pinger, err := probing.NewPinger(IPaddr)
	if err != nil {
		log.Fatalf("--- Erreur lors de la configuration du ping vers l'adresse spécifiée:\n%s", err)
	}
	pinger.Count = 1 // Nombre de paquets à envoyer.
	pinger.SetPrivileged(true)
	pinger.Timeout = time.Millisecond * 20 // Durée avant time out (en time.Duration).

	// Exécution du ping
	err = pinger.Run()
	if err != nil {
		log.Fatalf("--- Erreur lors de l'exécution du ping vers l'adresse spécifiée:\n%s", err)
	}

	// Résultat
	if pinger.Statistics().PacketsRecv == pinger.Statistics().PacketsSent {
		up = 1
	} else {
		up = 0
	}

	return up
}

func main() {
	fmt.Printf("aaa")
}
