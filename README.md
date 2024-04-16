# supervision-mikrotik-grafana

### Code et dépendances requises pour la supervision de routeurs Mikrotik via Grafana, Prometheus, SNMP Exporter.

# Mise en place

## Grafana

### Installation

On installe le package ```grafana``` :
```shell
sudo apt install grafana
```
puis on lance le service ```grafana-server```:
```shell
sudo systemctl start grafana-server.service
```