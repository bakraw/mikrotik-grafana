package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gorilla/mux"
	probing "github.com/prometheus-community/pro-bing"
)

type Router struct {
	IP     string  `json:"ip"`
	Lat    float64 `json:"lat"`
	Lon    float64 `json:"lon"`
	Status int     `json:"status"`
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
	var filePath string = strings.ReplaceAll(curDir, "api-json", "routers.json")

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

func getRoot(w http.ResponseWriter, r *http.Request) {

	json.NewEncoder(w).Encode(readJSON())

	fmt.Printf("Got GET request on /\n")
}

func handleRequests() {

	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", getRoot)

	log.Fatal(http.ListenAndServe("localhost:3333", router))
}

// Ping une adresse IP pour vérifier son état.
// Prend en entrée un adresse IP (string) et renvoie le statut (int, up = 1 et down = 0).
// On peut changer le nombre de paquets à envoyer et la durée avant time out.
func probeIP(IPaddr string) int {

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
		return 1
	} else {
		return 0
	}
}

func probeAll(routers []Router) {

	for {
		fmt.Println("--- Mise à jour du statut des routeurs...")

		for i := range routers {
			routers[i].Status = probeIP(routers[i].IP)
		}

		writeJSON(routers)
		fmt.Println("--- Mise à jour terminée.")
		time.Sleep(time.Second * 30)
	}

}

func main() {
	go probeAll(readJSON())
	handleRequests()
}
