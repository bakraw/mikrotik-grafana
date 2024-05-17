# mikrotik-grafana

### Code et fichiers de configuration pour la supervision de routeurs Mikrotik via Grafana, Prometheus, SNMP Exporter.

*src* contient le code source de *mikromap-cli* qui permet d'ajouter un nouveau routeur à tous les fichiers nécessaires, et de *mikromap-api*, un serveur HTTP qui transmet les informations au panel Geomap de Grafana.

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

### Configuration

#### Mikrotik

Sur les routeurs:
1. Supprimer la communauté *public* par défaut et en créer une nouvelle nommée *private*. 
2. Limiter la communauté à la lecture seule.
3. Mettre la sécurité sur *private*.
4. *SHA* en protocole d'authentification et *AES* en protocole de chiffrement.
5. Indiquer les mots de passe d'authentification et de chiffrement.
6. Limiter l'accès à la communauté à l'IP du serveur de supervision (*address* en CLI, champ *Addresses* dans WinBox).

Via WinBox, ces options se trouvent dans *IP* > *SNMP* puis *Communities*.

![Configuration communauté SNMP via WinBox](https://github.com/bakraw/mikrotik-grafana/assets/161661948/c6bbd279-0839-4abb-8183-d13a1e1e7255)

> N.B.- Pour l'instant, la paire de mots de passe doit être identique sur tous les routeurs / pare-feux. Le fichier de configuration de SNMP Exporter fait plus de 100 000 lignes de long et automatiser sa modification est une purge, tandis que Prometheus associe les configs d'authentification aux jobs et non aux cibles (donc il faudrait créer un job par routeur, et redémarrer Prometheus à chaque ajout). Désolé :(

#### Watchguard

Sur les pare-feux, veiller à ce que les protocoles SNMP (pour les infos) et ICMP (pour la vérification de l'état sur la carte) soit autorisé pour l'IP du serveur de supervision, et configurer la communauté de la même manière que sur les routeurs Mikrotik.

#### Fichier de configuration SNMP Exporter

Dans *snmp_config.yml*, indiquer les mot de passe d'authentification et de chiffrement SNMP (en haut du fichier, *password* et *priv_password*).

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
./grafana server --config=$HOME/mikrotik-grafana/conf/grafana_config.ini
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

> N.B.- Laisser les noms par défaut (sinon il faudra re-sélectionner les sources partout où elles sont utilisées).

![Config data source Prometheus](https://github.com/bakraw/mikrotik-grafana/assets/161661948/cd5f8abe-a194-4a92-9e77-a2ad1b673a86)

Dans la barre latérale: *Dashboards*, puis *New* > *Import* > *Upload dashboard JSON file* , et choisir ```~/mikrotik-grafana/conf/grafana_dashboard.json```. Séléctionner les sources de données précédemment crées.

Dans la barre latérale: *Administration* > *General* > *Default preferences*, sélectionner *Home Dashboard* = ```General/Supervision Mikrotik```.

Créer d'autres utilisateurs si besoin (dans la barre latérale: *Administration* > *Users and access* > *Users* puis *New user*)

> N.B.- Dans la configuration actuelle, tous les utilisateurs ont accès à tous les dashboards du dossier principal, donc éviter d'en créer d'autres. Sinon, les mettre dans des dossiers à accès restreint.

## Ajout et supression de routeur

### Ajout

L'ajout de routeur à la supervision se fait via *mikromap-cli*:
```bash
cd ~/mikrotik-grafana/bin/
./mikromap-cli
```

- Si l'adresse IP à ajouter correspond à un Watchguard, l'indiquer en ajoutant un *W* sans espace avant l'adresse IP pour éviter des problèmes de compatibilité (ex: ***W**8.8.8.8*)

> N.B.- L'adresse postale entrée n'a pas besoin d'être parfaitement écrite (pas besoin d'accents, tirets, etc.) mais veiller à inclure un minimum d'informations pour que l'API renvoie les bonnes coordonnées (ex: *1 rue leclerc st etienne* suffit à obtenir *1 Rue du Général Leclerc 42100 Saint-Étienne*)

> N.B.- Le nom d'utilisateur Grafana renseigné est comparé à celui renvoyé directement par Grafana, et doit donc **être identique** à celui du compte Grafana associé (pas grave si les majuscules sont différentes), sinon il n'apparaîtra pas sur le dashboard de cet utilisateur. Laisser le champ vide si le routeur ne doit être vu que par l'admin.

### Suppression

Pour supprimer un routeur, utiliser *mikromap-cli* et entrer un nombre négatif de routeurs à ajouter. 

> N.B.- Il n'y a besoin que de l'adresse IP pour supprimer un routeur, et le préfixe *W* n'est pas nécessaire pour désigner un Watchguard.
