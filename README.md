# mikrotik-grafana

### Code et dépendances requises pour la supervision de routeurs Mikrotik via Grafana, Prometheus, SNMP Exporter.

```geomap-routeurs``` contient ```outil-cli``` qui permet d'ajouter un nouveau routeur à tous les fichiers nécessaires, et ```api-json```, un serveur HTTP qui transmet les informations au panel geomap de Grafana.

```fichiers-config``` contient les fichiers de configuration utilisés par les composants.

## Mise en place

### Téléchargement

Télécharger les binaires stand-alone de [Grafana](https://grafana.com/get/?tab=self-managed), [Prometheus](https://prometheus.io/download/), [SNMP-Exporter](https://github.com/prometheus/snmp_exporter/releases) et de [ce dépôt](https://github.com/bakraw/mikrotik-grafana/releases).

Depuis le dossier où ils ont été téléchargés, placer dans le répertoire personnel et extraire:
```bash
mv *.tar.gz ~
cd ~
tar -xf *.tar.gz
rm -rf *.tar.gz
```

### Lancement

Lancer Prometheus:
```bash
~/prometheus*/prometheus --config.file=$HOME/mikrotik-grafana-release/fichiers-config/prometheus_config.yml
```

Lancer SNMP Exporter:
```bash
~/snmp_exporter*/snmp_exporter --config.file=$HOME/mikrotik-grafana-release/fichiers-config/snmp_config.yml
```

Lancer Grafana:
```bash
cd ~/grafana*/bin/
./grafana server
```

### Grafana

Ouvrir l'interface web à l'adresse ```localhost:3000```.

Se connecter (UN:```admin``` / PW:```admin```).

Dans la barre latérale: *Administration* > *Plugins and data* > *Plugins*
A côté de la barre de recherche: *State* = ```All```

![Menu Plugin de Grafana](https://github.com/bakraw/mikrotik-grafana/assets/161661948/ee092fb0-bfa8-4260-801c-b95fcdd0b77b)

Installer le plugin *JSON API* de Marcus Olsson

![JSON API](https://github.com/bakraw/mikrotik-grafana/assets/161661948/28660e68-0f56-4d53-92a4-50dd030e6fb7)