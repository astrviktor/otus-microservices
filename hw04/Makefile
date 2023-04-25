# helm install

helm_install_prometheus_stack:
	@helm install prometheus-stack prometheus-community/kube-prometheus-stack \
      --namespace prometheus --create-namespace \
      --set prometheus.prometheusSpec.podMonitorSelectorNilUsesHelmValues=false \
      --set prometheus.prometheusSpec.serviceMonitorSelectorNilUsesHelmValues=false

helm_install_ingress:
	@helm install ingress-nginx ingress-nginx/ingress-nginx \
      --namespace ingress-nginx --create-namespace \
      --set controller.metrics.serviceMonitor.enabled=true \
      --set controller.metrics.enabled=true \
      --set controller.metrics.serviceMonitor.additionalLabels.release="prometheus" \
      -f ./helm-charts/ingress/nginx-ingress.yaml

helm_install_service_and_postgres:
	@helm install storage ./helm-charts/postgres/postgresql-12.2.8.tgz \
       --set auth.postgresPassword=password \
       --set auth.database=users
	@sleep 60
	@helm install crudservice ./helm-charts/crudservice/

helm_install_postgres_exporter:
	@helm install exporter prometheus-community/prometheus-postgres-exporter \
       -f ./helm-charts/exporter/values.yaml


# helm uninstall

helm_uninstall_postgres_exporter:
	@helm uninstall exporter

helm_uninstall_service_and_postgres:
	@helm uninstall crudservice
	@helm uninstall storage

helm_uninstall_ingress:
	@helm uninstall ingress-nginx

helm_uninstall_prometheus_stack:
	@helm uninstall prometheus-stack


# newman test postman collection

newman_run:
	@newman run ./postman/crudservice.postman_collection.json