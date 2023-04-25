## Выполнение домашнего задания

### В прошлых ДЗ
- Установка minikube
- Установка helm

### Запуск minikube без Ingress
```
Для того, чтобы снимать метрики с ingress, лучше не использовать плагин в minikube 
а поставить ingress отдельно

# Запуск minikube
minikube start --cpus=2 --cni=flannel --install-addons=true --kubernetes-version=stable --memory=8g
```

### Установка Prometheus и Grafana

```
# добавление репозиториев

helm repo add prometheus-community https://prometheus-community.github.io/helm-charts
helm repo update

# установка

helm upgrade --install prometheus-stack prometheus-community/kube-prometheus-stack \
  --namespace prometheus --create-namespace \
  --set prometheus.prometheusSpec.podMonitorSelectorNilUsesHelmValues=false \
  --set prometheus.prometheusSpec.serviceMonitorSelectorNilUsesHelmValues=false

или

make helm_install_prometheus_stack

# проброс портов Prometheus

kubectl port-forward service/prometheus-operated 9090 --namespace=prometheus

# проброс портов Grafana

kubectl port-forward service/prometheus-stack-grafana 3000:80 --namespace=prometheus

# доступ в grafana
admin / prom-operator
```


### Установка Ingress
```
Нужна версия Ingress с метриками, поэтому нужно ставить отдельно

helm upgrade --install ingress-nginx ingress-nginx/ingress-nginx \
  --namespace ingress-nginx --create-namespace \
  --set controller.metrics.serviceMonitor.enabled=true \
  --set controller.metrics.enabled=true \
  --set controller.metrics.serviceMonitor.additionalLabels.release="prometheus" \
  -f ./helm-charts/ingress/nginx-ingress.yaml

или

make helm_install_ingress

# проверка и добавление адреса в /etc/hosts (после деплоя сервиса)

kubectl get ingress

NAME                  CLASS   HOSTS           ADDRESS          PORTS   AGE
crudservice-ingress   nginx   arch.homework   10.111.249.253   80      6m34s

sudo nano /etc/hosts
```


### Установка Сервиса и Postgres

```
helm install storage ./helm-charts/postgres/postgresql-12.2.8.tgz \
  --set auth.postgresPassword=password \
  --set auth.database=users

helm install crudservice ./helm-charts/crudservice/

или

make helm_install_service_and_postgres
```

### Установка prometheus-postgres-exporter
```
helm install exporter prometheus-community/prometheus-postgres-exporter \
  -f ./helm-charts/exporter/values.yaml

или

make helm_install_postgres_exporter
```


### Проверка, что метрики попадают в Prometheus
```
# проброс порта
kubectl port-forward service/prometheus-operated 9090 --namespace=prometheus

по адресу
http://127.0.0.1:9090/targets?search=

должны быть следующие таргеты:

serviceMonitor/default/crudservice-monitor/0 (2/2 up)
serviceMonitor/default/exporter-prometheus-postgres-exporter/0 (1/1 up)
serviceMonitor/ingress-nginx/ingress-nginx-controller/0 (1/1 up)

```


### Нагрузка на сервис
```
Для тестирования грфиков сделан специальный метод /testing
в котором есть разная задержка и разные коды ответов

Для нагрузки будет использоваться Apache Bench

sudo apt-get install apache2-utils 
ab -V

# нагрузка 600 запросов, 10 запросов в секунду
ab -n 600 -c 10 http://arch.homework/testing
```



### Настройка дашбордов в Grafana
```
# проброс портов Grafana
kubectl port-forward service/prometheus-stack-grafana 3000:80 --namespace=prometheus

# доступ в grafana
http://127.0.0.1:3000/login

admin / prom-operator

# Готовые дашборды (импорт по номеру)
http://127.0.0.1:3000/dashboard/import

Postgres - 455 или 9628
Node Exporter - 1860
Kubernetes Cluster - 6417
NGINX Ingress controller - 9614

# Метрики сервиса
RPS
sum(rate(response_duration_bucket[1m])) by (method, path, code)  

Error Rate
sum(rate(response_duration_bucket{code=~"5.."}[1m])) by (method, path, code) 
 
Latency (average response time)
sum(rate(response_duration_sum[1m]) / rate(response_duration_count[1m])) by (path, method, code, pod)

Latency (0.50, 0.95, 0.99, 1.00)
histogram_quantile(0.50, sum(rate(response_duration_bucket[1m])) by (le, path))
histogram_quantile(0.95, sum(rate(response_duration_bucket[1m])) by (le, path))
histogram_quantile(0.99, sum(rate(response_duration_bucket[1m])) by (le, path))
histogram_quantile(1.00, sum(rate(response_duration_bucket[1m])) by (le, path))

Ingress Request Volume
round(sum(irate(nginx_ingress_controller_requests{}[1m])) by (ingress), 0.001)

Ingress Error Rate (5xx responses). %
sum(rate(nginx_ingress_controller_requests{status=~"5.*"}[2m])*100) by (ingress) / sum(rate(nginx_ingress_controller_requests{}[2m])) by (ingress)

RAM
go_memstats_alloc_bytes{service="crudservice"}

Cpu usage
rate(process_cpu_seconds_total{service="crudservice"}[1m]) * 1000

json-дашборд ./dashboards/crudservice.json
```

### Удаление сервисов через helm или Makefile
```
helm uninstall exporter
helm uninstall crudservice
helm uninstall storage
helm uninstall ingress-nginx
helm uninstall prometheus-stack

или

make helm_uninstall_postgres_exporter
make helm_uninstall_service_and_postgres
make helm_uninstall_ingress
make helm_uninstall_prometheus_stack
```

### Остановка и удаление minikube (если нужно)
```
minikube stop
minikube delete
```

### Скриншоты
![Alt text](./pictures/postgres1.jpg?raw=true "")
![Alt text](./pictures/postgres2.jpg?raw=true "")

![Alt text](./pictures/kubernetes.jpg?raw=true "")

![Alt text](./pictures/nodeexporter.jpg?raw=true "")

![Alt text](./pictures/nginx.jpg?raw=true "")

![Alt text](./pictures/crudservice.jpg?raw=true "")

