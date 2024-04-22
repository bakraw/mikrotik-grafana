// Dernière mise à jour: avril 2024
package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
)

type Router struct {
	IP     string  `json:"ip"`
	Lat    float64 `json:"lat"`
	Lon    float64 `json:"lon"`
	Status int     `json:"status"`
}

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
		log.Fatalf("--- Erreur lors du traitement des données JSON reçues:\n%s", err)
	}

	// Récupération des coordonnées
	lat, lon = target.Features[0].Geometry.Coordinates[1], target.Features[0].Geometry.Coordinates[0]

	return lat, lon
}

// Renvoie le chemin vers le fichier JSON.
// Ne prend rien en entrée et renvoie le chemin (string).
// A modifier si besoin de mettre le fichier ailleurs que dans le répertoire parent.
func getPath() string {

	// Récupération du chemin vers le fichier
	curDir, err := os.Getwd()
	if err != nil {
		log.Fatalf("--- Erreur lors de la récupération du répertoire courant:\n%s", err)
	}
	var filePath string = strings.ReplaceAll(curDir, "outil-cli", "routers.json")

	return filePath
}

// Récupère les données du fichier de stockage JSON.
// Ne prend rien en entrée et renvoie les données dans un struct []Router.
func readJSON() []Router {

	var data []Router

	// Lecture du fichier
	content, err := os.ReadFile(getPath())
	if err != nil {
		log.Fatalf("--- Erreur lors de la lecture du fichier JSON:\n%s", err)
	}

	// Traitement des données
	err = json.NewDecoder(bytes.NewBuffer(content)).Decode(&data)
	if err != nil {
		log.Fatalf("--- Erreur lors du traitement des données du fichier JSON:\n%s", err)
	}

	return data
}

// Ecrit par-dessus le fichier JSON.
// Prend en entrée les données à écrire et ne renvoie rien.
func writeJSON(data []Router) {

	content, err := os.OpenFile(getPath(), os.O_WRONLY, os.ModePerm)
	if err != nil {
		log.Fatalf("--- Erreur lors de l'ouverture du fichier JSON pour écriture:\n%s", err)
	}

	err = json.NewEncoder(content).Encode(data)
	if err != nil {
		log.Fatalf("--- Erreur lors de l'écriture du fichier JSON:\n%s", err)
	}
}

// Pour l'instant, ne sert qu'à tester.
// Servira probablement de menu.
func main() {

	var addrPost string
	var addrIP string

	// Récupération adresse postale
	fmt.Print("Adresse postale >> ")
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
	fmt.Print("Adresse IP >> ")
	if scanner.Scan() {
		addrIP = scanner.Text()
	}
	fmt.Println(addrIP)

	data := readJSON()

	// Ajout d'un nouveau routeur
	newRouter := Router{
		IP:     addrIP,
		Lat:    lat,
		Lon:    lon,
		Status: 0,
	}

	data = append(data, newRouter)

	writeJSON(data)

	// Récupération données du fichier
	data = readJSON()

	fmt.Printf("%+v", data)
}
