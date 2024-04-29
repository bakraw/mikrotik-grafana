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

// Structure routers.json
type Router struct {
	IP      string  `json:"ip"`
	Lat     float64 `json:"lat"`
	Lon     float64 `json:"lon"`
	Adresse string  `json:"adresse"`
	Statut  int     `json:"statut"`
}

// Structure prometheus_targets.json
type PromTargets struct {
	Labels  Labels   `json:"labels"`
	Targets []string `json:"targets"`
}
type Labels struct {
	Job string `json:"job"`
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
// Prend en entrée des données JSON ([]byte) et renvoie latitude (float64), longitude (float64) et adresse (string).
// Doit être adaptée selon l'API.
func extractCoords(data []byte) (float64, float64, string) {

	var lat, lon float64
	var adresse string

	// Structure de la réponse JSON.
	// Doit être adaptée selon l'API. On peut omettre les champs inutiles.
	// Des outils existent pour le générer automatiquement (ex: https://mholt.github.io/json-to-go/)
	type Geometry struct {
		Coordinates []float64 `json:"coordinates"`
	}
	type Properties struct {
		Label string `json:"label"`
	}
	type Features struct {
		Geometry   Geometry   `json:"geometry"`
		Properties Properties `json:"properties"`
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
	lat, lon, adresse = target.Features[0].Geometry.Coordinates[1], target.Features[0].Geometry.Coordinates[0], target.Features[0].Properties.Label

	return lat, lon, adresse
}

// Renvoie le chemin vers le fichier JSON spécifié.
// Prend le nom du fichier en entrée et renvoie le chemin (string).
// A modifier si besoin de mettre le fichier ailleurs.
func getPath(target string) string {

	filePath := fmt.Sprintf("%s/mikrotik-grafana/conf/%s", os.Getenv("HOME"), target)
	return filePath
}

// Récupère les données du fichier de stockage JSON.
// Ne prend rien en entrée et renvoie les données dans un slice de struct []Router.
func readJSON() []Router {

	var data []Router

	// Lecture du fichier
	content, err := os.ReadFile(getPath("routers.json"))
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

	// Ouverture du fichier
	content, err := os.OpenFile(getPath("routers.json"), os.O_WRONLY|os.O_TRUNC, os.ModePerm)
	if err != nil {
		log.Fatalf("--- Erreur lors de l'ouverture du fichier JSON pour écriture:\n%s", err)
	}

	// Ecriture du fichier
	enc := json.NewEncoder(content)
	enc.SetIndent("", "    ")
	err = enc.Encode(data)
	if err != nil {
		log.Fatalf("--- Erreur lors de l'écriture du fichier JSON:\n%s", err)
	}
}

// Récupère les données du fichier des cibles Prometheus JSON.
// Ne prend rien en entrée et renvoie les données dans un slice de struct []PromTargets.
func readPromTargets() []PromTargets {

	var data []PromTargets

	// Lecture du fichier
	content, err := os.ReadFile(getPath("prometheus_targets.json"))
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

// Ecrit par-dessus le fichier de cibles Prometheus JSON.
// Prend en entrée les données à écrire et ne renvoie rien.
func writePromTargets(data []PromTargets) {

	// Ouverture du fichier
	content, err := os.OpenFile(getPath("prometheus_targets.json"), os.O_WRONLY|os.O_TRUNC, os.ModePerm)
	if err != nil {
		log.Fatalf("--- Erreur lors de l'ouverture du fichier JSON pour écriture:\n%s", err)
	}

	// Ecriture du fichier
	enc := json.NewEncoder(content)
	enc.SetIndent("", "    ")
	err = enc.Encode(data)
	if err != nil {
		log.Fatalf("--- Erreur lors de l'écriture du fichier JSON:\n%s", err)
	}
}

// Fonction principale qui ajoute un routeur aux fichiers.
// Ne prend rien en entrée et ne renvoie rien.
func addRouter() {

	var addrPost string
	var addrIP string

	fmt.Println("--- Ajouter un routeur à la supervision")

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
	lat, lon, adresse := extractCoords(resBody)
	fmt.Printf("- %s\n- %f, %f\n", adresse, lat, lon)

	// Récupération adresse IP
	fmt.Print("Adresse IP >> ")
	_, err := fmt.Scanln(&addrIP)
	if err != nil {
		log.Fatalf("--- Erreur lors de la récupération de la saisie:\n%s", err)
	}

	// Ajout d'un nouveau routeur dans routers.json
	dataR := readJSON()

	newRouter := Router{
		IP:      addrIP,
		Lat:     lat,
		Lon:     lon,
		Adresse: adresse,
		Statut:  0,
	}

	dataR = append(dataR, newRouter)
	writeJSON(dataR)

	// Ajout IP dans prometheus_targets.json
	dataT := readPromTargets()
	dataT[0].Targets = append(dataT[0].Targets, addrIP)
	writePromTargets(dataT)

	fmt.Println("--- Routeur ajouté")
}

// Fonction principale qui retire un routeur des fichiers.
// Ne prend rien en entrée et ne renvoie rien.
func removeRouter() {
	var addrIP string

	fmt.Println("--- Retirer un routeur de la supervision")

	fmt.Print("Adresse IP du routeur à supprimer >>> ")
	_, err := fmt.Scanln(&addrIP)
	if err != nil {
		log.Fatalf("--- Erreur lors de la récupération de la saisie:\n%s", err)
	}

	dataR := readJSON()
	dataT := readPromTargets()

	for i, v := range dataR {
		if v.IP == addrIP {
			dataR = append(dataR[0:i], dataR[i+1:]...)
		}
	}
	writeJSON(dataR)

	for i, v := range dataT[0].Targets {
		if v == addrIP {
			dataT[0].Targets = append(dataT[0].Targets[0:i], dataT[0].Targets[i+1:]...)
		}
	}
	writePromTargets(dataT)

	fmt.Println("--- Routeur supprimé")
}

func main() {

	var n int
	fmt.Print("Nombre de routeurs à ajouter >>> ")
	_, err := fmt.Scanln(&n)
	if err != nil {
		log.Fatalf("--- Erreur lors de la récupération de la saisie:\n%s", err)
	}

	if n >= 0 {
		for i := 0; i < n; i++ {
			addRouter()
		}
	} else {
		for i := 0; i < -n; i++ {
			removeRouter()
		}
	}
}
