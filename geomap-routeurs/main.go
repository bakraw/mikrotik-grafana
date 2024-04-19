// Dernière mise à jour: avril 2024
package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	probing "github.com/prometheus-community/pro-bing"
)

// Fait un appel à l'API renseignée.
// Prend en entrée une adresse (string), renvoie le code de statut (int) et le corps ([]byte) de la réponse.
// L'API Adresse du gouvernement est gratuite et fonctionne parfaitement pour la France (50 calls/IP/sec).
// Si besoin d'utiliser une autre API, penser à changer la partie formatage et la fonction extractCoords.
func geoAPI(addr string) ([]byte, int) {

	// Formatage de la requête
	var addrConcat string = strings.ReplaceAll(addr, " ", "+")
	var reqURL string = fmt.Sprintf("https://api-adresse.data.gouv.fr/search/?q=%s&limit=1", addrConcat) // Modifier en fonction de l'API à utiliser

	// Exécution de la requête
	res, err := http.Get(reqURL)
	if err != nil {
		log.Fatalf("--- Erreur lors de l'appel à l'API de géocodage:\n%s", err)
	}

	// Traitement de la réponse
	resData, err := io.ReadAll(res.Body)
	if err != nil {
		log.Fatalf("--- Erreur lors du traitement de la réponse de l'API de géocodage:\n%s", err)
	}

	return resData, res.StatusCode
}

// Traite le JSON reçu et en extrait les coordonnées de l'adresse.
// Prend en entrée des données JSON ([]byte) et renvoie latitude (float64) et longitude (float64).
// Doit être adaptée selon l'API.
func extractCoords(data []byte) (float64, float64) {

	var lat, lon float64

	// Structure de la réponse JSON.
	// Doit être adaptée selon l'API. On peut omettre les champs inutiles.
	// Des outils existent pour le générer automatiquement (ex: https://mholt.github.io/json-to-go/)
	type Geometry struct {
		Coordinates []float64 `json:"coordinates"`
	}
	type Features struct {
		Geometry Geometry `json:"geometry"`
	}
	type Data struct {
		Features []Features `json:"features"`
	}

	var target Data

	// Traitement des données JSON
	err := json.Unmarshal(data, &target)
	if err != nil {
		log.Fatalf("--- Erreur lors du traitement des données JSON:\n%s", err)
	}

	// Récupération des coordonnées
	lat, lon = target.Features[0].Geometry.Coordinates[1], target.Features[0].Geometry.Coordinates[0]

	return lat, lon
}

// Ping une adresse IP pour vérifier son état.ody
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

// Pour l'instant, ne sert qu'à tester.
// Servira probablement de menu.
func main() {
	var addrPost string
	var addrIP string

	// Récupération adresse postale
	fmt.Print("Adresse >> ")
	scanner := bufio.NewScanner(os.Stdin)
	if scanner.Scan() {
		addrPost = scanner.Text()
	}

	// Récupération coordonnées géographiques
	resBody, resCode := geoAPI(addrPost)
	if resCode != 200 {
		log.Fatalf("--- Erreur lors de l'appel à l'API de géocodage (code %d)", resCode)
	}
	lat, lon := extractCoords(resBody)
	fmt.Printf("%f %f\n", lat, lon)

	// Récupération adresse IP
	fmt.Print("IP >> ")
	if scanner.Scan() {
		addrIP = scanner.Text()
	}

	// Traitement adresse IP
	status := probeIP(addrIP)
	fmt.Println(status)
}
