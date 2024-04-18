# supervision-mikrotik-grafana

### Code et dépendances requises pour la supervision de routeurs Mikrotik via Grafana, Prometheus, SNMP Exporter.

Le dossier ```geomap-routeurs``` contiendra tout ce qui sera nécessaire à l'affichage des routeurs et de leur statut sur un panneau Geomap de Grafana.

Cela incluera une application principale qui permettra de manipuler un stockage (JSON / BDD ?), afin d'y mettre l'IP, les coordonnées GPS et le statut de chaque routeur. L'application devra aussi pouvoir retrouver les coordonnées à partir d'une adresse (géocodage), et mettre à jour la liste des IP cibles de SNMP Exporter.

D'autres dossiers serviront sans doute à contenir les binaires nécessaires au projet complet (Prometheus, Grafana, SNMP Exporter, etc.), l'idée étant de conteneuriser le tout.