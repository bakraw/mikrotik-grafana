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

// Structure routers.json
type Router struct {
	IP       string  `json:"ip"`
	Lat      float64 `json:"lat"`
	Lon      float64 `json:"lon"`
	Adresse  string  `json:"adresse"`
	Username string  `json:"username"`
	Statut   int     `json:"statut"`
}

// Renvoie le chemin vers le fichier JSON.
// Ne prend rien en entrée et renvoie le chemin (string).
// A modifier si besoin de mettre le fichier ailleurs que dans le répertoire parent.
func getPath() string {

	filePath := fmt.Sprintf("/home/%s/mikrotik-grafana/conf/routers.json", os.Getenv("SUDO_USER"))
	return filePath
}

// Récupère les données du fichier de stockage JSON.
// Ne prend rien en entrée et renvoie les données dans un struct []Router.
func readJSON() []Router {

	var data []Router

	// Lecture du fichier
	content, err := os.ReadFile(getPath())
	if err != nil {
		log.Fatalf("--- Erreur lors de la lecture du fichier JSON (vérifier que l'exécutable a bien été lancé en sudo):\n%s", err)
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
	content, err := os.OpenFile(getPath(), os.O_WRONLY, os.ModePerm)
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

// Traite les requêtes HTTP GET.
// Renvoie le contenu de routers.json qui concerne l'utilisateur Grafana qui fait le call.
// Prend en entrée un http.responseWriter et une http.Request.
// Ne devrait être appelée que via HandleFunc().
func getMikromap(w http.ResponseWriter, r *http.Request) {

	// Récupération du nom d'utilisateur transmis par Grafana
	user := r.URL.Query().Get("user")

	fmt.Printf("Requête GET entrante sur /mikromap (user = %s)\n", user)

	// Récupération routers.json dans un struct
	dataRouters := readJSON()

	// Suppression dans le struct de tous les routeurs dont le Username n'est pas identique au paramètre reçu.
	// Si le paramètre vaut admin, on saute cette étape.
	if user != "admin" {
		// On parcourt le slice dans le sens inverse pour ne pas modifier des éléments pas encore parcourus.
		for i := len(dataRouters) - 1; i >= 0; i-- {
			v := dataRouters[i]
			if !(strings.EqualFold(v.Username, user)) {
				dataRouters = append(dataRouters[0:i], dataRouters[i+1:]...)
			}
		}
	}

	// Envoi du struct modifié
	json.NewEncoder(w).Encode(dataRouters)
}

// Traite les requêtes HTTP entrantes.
// Ne prend rien en entrée et en renvoie rien.
// Fonction sans condition de sortie.
func handleRequests() {

	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/mikromap", getMikromap)

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
	pinger.Timeout = time.Millisecond * 500 // Durée avant time out (en time.Duration).

	// Exécution du ping
	err = pinger.Run()
	if err != nil {
		log.Fatalf("--- Erreur lors de l'exécution du ping vers l'adresse spécifiée:\n%s", err)
	}

	// Résultat
	if pinger.Statistics().PacketsRecv == pinger.Statistics().PacketsSent {
		return 1
	}
	return 0
}

// Teste toutes les IPs mentionnées dans un slice de struct puis ré-écrit le fichier JSON.
// Ne prend rien en entrée et ne renvoie rien.
// Fonction sans condition de sortie.
func probeAll() {
	var routers []Router

	for {
		fmt.Println("--- Mise à jour du statut des routeurs...")

		routers = readJSON()

		// Test des IPs
		for i := range routers {
			routers[i].Statut = probeIP(routers[i].IP)
		}

		// Ecriture du fichier JSON
		writeJSON(routers)

		fmt.Println("--- Mise à jour terminée.")
		time.Sleep(time.Second * 30) // Durée entre chaque rafraîchissement
	}
}

func main() {
	go probeAll() // Goroutine de test des IPs en parallèle du traitement des requêtes HTTP.
	handleRequests()
}
