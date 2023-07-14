## Выполнение домашнего задания

Реализованы:
- сервис заказов `orderservice`
- сервис платежей `payment`
- сервис склада `warehouse`
- сервис доставки `delivery`

По сервисам платежей, склада, доставки частичная реализация, 
для проверки распределенной транзакции

Используется паттерн распределенной транзакции - сага (оркестрация) 

Используется следующий алгоритм:
1) Запрос на создание заказа приходит в сервис заказов `orderservice`
2) Сервис заказов `orderservice` отправляет запросы в связанные сервисы 
- платежей `payment`
- склада `warehouse`
- доставки `delivery`
3) Если все запросы прошли успешно, создается заказ в сервисе заказов `orderservice`
4) Если хотя бы один из запросов прошел неуспешно, заказ в сервисе заказов `orderservice` не создается,
при необходимости отправляются компенсирующие запросы в связанные сервисы 


### Установка сервисов и postgres (namespace default)

```
helm install storage ./deployments/kubernetes/helm-charts/postgres/ \
    --set auth.postgresPassword=password \
    --set auth.database=orders

helm install orderservice ./deployments/kubernetes/helm-charts/orderservice/
helm install payment ./deployments/kubernetes/helm-charts/payment/
helm install warehouse ./deployments/kubernetes/helm-charts/warehouse/
helm install delivery ./deployments/kubernetes/helm-charts/delivery/

или

make helm_install
```

### Проверка тестового сценария Postman (newman)
```
newman run ./postman/saga.postman_collection.json --verbose

или

make newman_run
```

### Удаление сервисов и postgres
```
helm uninstall orderservice
helm uninstall payment
helm uninstall warehouse
helm uninstall delivery
helm uninstall storage

или

make helm_uninstall
```
