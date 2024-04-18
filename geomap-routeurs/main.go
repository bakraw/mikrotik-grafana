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
)

// Fait un appel à l'API renseignée.
// Prend en entrée une adresse (string), renvoie le code de statut (int) et le corps ([]byte) de la réponse.
// L'API Adresse du gouvernement est gratuite et fonctionne parfaitement pour la France (50 calls/IP/sec).
// Si besoin d'utiliser une autre API, penser à changer la partie formatage et la fonction extractCoords.
func geoAPI(addr string) (int, []byte) {

	// Formatage de la requête
	var addrConcat string = strings.ReplaceAll(addr, " ", "+")
	fmt.Println(addrConcat)
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

	return res.StatusCode, resData
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

// Pour l'instant, ne sert qu'à tester.
// Servira probablement de menu.
func main() {
	var input string

	fmt.Print("Adresse >> ")
	scanner := bufio.NewScanner(os.Stdin)
	if scanner.Scan() {
		input = scanner.Text()
	}

	fmt.Println(input)
	resCode, resBody := geoAPI(input)
	fmt.Printf("%d\n%s\n", resCode, resBody)
	lat, lon := extractCoords(resBody)
	fmt.Print(lat, lon)
}
