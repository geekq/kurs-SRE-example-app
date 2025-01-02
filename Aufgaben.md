## Prometheus einrichten

* [x] Server mit einem privaten Netzwerk erstellen.
* [x] DNS Eintrag anpassen.

```
ssh root@157.90.119.180 # Ist nur ein Beispiel. Bitte die IP Adresse anpassen.
hostname ........de
apt-get update
apt-get install prometheus
apt-get install prometheus-alertmanager
# prometheus-nginx-exporter
```

Um die Alerts per Email zu bekommen:

```
apt-get install mailutils
mail
```

### Alternative Anleitungen

Komplizierter, zeigt aber die einzelnen Schritte und kann als Basis für eigene
Automatisierung dienen. Nutzt die neuesten Binaries aus dem Prometheus Release
https://www.cherryservers.com/blog/install-prometheus-ubuntu

Hier bereits mit
[ansible](https://prometheus-community.github.io/ansible/branch/main/prometheus_role.html#ansible-collections-prometheus-prometheus-prometheus-role)
automatisiert.


## Metriken erkunden

Metrik-Typen:

* counter
* gauge
* histogram, besonders für Timing

https://prometheus.io/docs/concepts/metric_types/

[Abfragesprache](https://prometheus.io/docs/prometheus/latest/querying/basics/)

```
shop_queue_length
node_cpu_seconds_total

shop_request_duration_seconds_bucket
shop_request_duration_seconds_bucket{endpoint="/metrics"}
shop_request_duration_seconds_bucket{endpoint="/metrics"}[2m]
rate(shop_request_duration_seconds_bucket{endpoint="/metrics"}[2m])
sum(rate(shop_request_duration_seconds_bucket{endpoint="/metrics", le="0.5"}[2m]))
histogram_quantile(0.9, rate(shop_request_duration_seconds_bucket{endpoint="/metrics"}[2m]))

sum(rate(shop_request_duration_seconds_bucket{endpoint="/metrics", le="0.5"}[2m])) / sum(rate(shop_request_duration_seconds_bucket{endpoint="/metrics", le="+Inf"}[2m]))
```

