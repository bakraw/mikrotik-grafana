# mikrotik-grafana

### Code et fichiers de configuration pour la supervision de routeurs Mikrotik via Grafana, Prometheus, SNMP Exporter.

*src* contient *mikromap-cli* qui permet d'ajouter un nouveau routeur à tous les fichiers nécessaires, et *mikromap-api*, un serveur HTTP qui transmet les informations au panel Geomap de Grafana.

*conf* contient les fichiers de configuration utilisés par les composants.

Le dashboard est une version modifiée de [Mikrotik monitoring](https://grafana.com/grafana/dashboards/14420-mikrotik-monitoring/) par *igorkha*.

## Mise en place

### Téléchargement

Télécharger les binaires stand-alone de [Grafana](https://grafana.com/get/?tab=self-managed), [Prometheus](https://prometheus.io/download/), [SNMP-Exporter](https://github.com/prometheus/snmp_exporter/releases) et de [ce dépôt](https://github.com/bakraw/mikrotik-grafana/releases) (ou le clôner et le build).

Depuis le dossier où ils ont été téléchargés, placer dans le répertoire personnel et extraire:
```bash
mv *.tar.gz ~
cd ~
tar -xf mikrotik*
tar -xf grafana*
tar -xf prometheus*
tar -xf snmp_exporter*
rm -rf *.tar.gz
```

### Lancement

> N. B.- Ajouter des services *systemd* pour chaque exécutable est recommandé pour éviter d'avoir à les relancer manuellement à chaque redémarrage.

Lancer Prometheus:
```bash
~/prometheus*/prometheus --config.file=$HOME/mikrotik-grafana/conf/prometheus_config.yml
```

Lancer SNMP Exporter:
```bash
~/snmp_exporter*/snmp_exporter --config.file=$HOME/mikrotik-grafana/conf/snmp_config.yml
```

Lancer Grafana:
```bash
cd ~/grafana*/bin/
./grafana server
```

Lancer l'API pour la carte:
```bash
cd ~/mikrotik-grafana/bin/
sudo ./mikromap-api
```

> N. B.- L'API doit obligatoirement être lancée en sudo pour que les pings fonctionnent.

### Grafana

Ouvrir l'interface web à l'adresse ```localhost:3000```.

Se connecter (username:```admin``` / password:```admin```).

Dans la barre latérale: *Administration* > *Plugins and data* > *Plugins* et a côté de la barre de recherche: *State* = ```All```

![Menu Plugin de Grafana](https://github.com/bakraw/mikrotik-grafana/assets/161661948/ee092fb0-bfa8-4260-801c-b95fcdd0b77b)

Installer les plugins *JSON API* de Marcus Olsson, et *Orchestra Cities Map* de Orchestra Cities by Martel

![JSON API](https://github.com/bakraw/mikrotik-grafana/assets/161661948/28660e68-0f56-4d53-92a4-50dd030e6fb7)

Dans la barre latérale: *Connections* > *Data sources*

Ajouter deux sources de données:
1. Prometheus (*Prometheus server URL* = ```http://localhost:9090```)
2. JSON API (*URL* = ```http://localhost:3333```)

> Laisser les noms par défaut (sinon il faudra re-sélectionner les sources partout où elles sont utilisées).

![Config data source Prometheus](https://github.com/bakraw/mikrotik-grafana/assets/161661948/cd5f8abe-a194-4a92-9e77-a2ad1b673a86)

Dans la barre latérale: *Dashboards*, puis *New* > *Import* > *Upload dashboard JSON file* , et choisir ```~/mikrotik-grafana/conf/grafana_dashboard.json```. Séléctionner les sources de données précédemment crées.

## Ajout et supression de routeur

L'ajout de routeur à la supervision se fait via *mikromap-cli*:
```bash
cd ~/mikrotik-grafana/bin/
./mikromap-cli
```

> N. B.- L'adresse entrée n'a pas besoin d'être parfaitement écrite (pas besoin d'accents, tirets, etc.) mais veiller à inclure un minimum d'informations pour que l'API renvoie les bonnes coordonnées (ex: *1 rue leclerc st etienne* suffit à obtenir *1 Rue du Général Leclerc 42100 Saint-Étienne*)

### Supression

Pour supprimer un routeur utiliser *mikromap-cli* et entrer un nombre négatif de routeurs à ajouter. Il n'y a besoin que de l'adresse IP pour supprimer un routeur.