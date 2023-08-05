## Выполнение домашнего задания

### Схемы вхаимодействия

http взаимодействие:

![Alt text](./pictures/http.jpg?raw=true "")

http взаимодействие с использованием брокера сообщений для нотификаций:

![Alt text](./pictures/http-event.jpg?raw=true "")

event collaboration взаимодействие с использованием брокера сообщений:
![Alt text](./pictures/http-event-collaboration.jpg?raw=true "")


Для реализации выбрана схема 

- http взаимодействия с использованием брокера сообщений для нотификаций

Реализованы:
- сервис заказов `orderservice`
- сервис биллинга `billing`
- сервис нотификаций `notification`


Используется следующий алгоритм:
1) Создается клиент и происходит пополнение баланса в сервисе `billing`
2) Создается заказ на определенную сумму в сервисе заказов `orderservice`
3) Если для заказа достаточно средств:
- происходит списание средств в сервисе `billing`
- происходит отправка успешного уведомления в сервис `notification`
4) Если для заказа недостаточно средств:
- происходит ошибка списания средств в сервисе `billing`
- происходит отправка уведомления с ошибкой в сервис `notification`


### Установка сервисов и postgres (namespace default)

```
helm install storage ./deployments/kubernetes/helm-charts/postgres/ \
    --set auth.postgresPassword=password \
    --set auth.database=orders

helm install rabbitmq ./deployments/kubernetes/helm-charts/rabbitmq/

helm install orderservice ./deployments/kubernetes/helm-charts/orderservice/
helm install billing ./deployments/kubernetes/helm-charts/billing/
helm install notification ./deployments/kubernetes/helm-charts/notification/

или

make helm_install
```

### Проверка тестового сценария Postman (newman)
```
newman run ./postman/hw10.postman_collection.json --verbose

или

make newman_run
```

### Удаление сервисов и postgres
```
helm uninstall orderservice
helm uninstall billing
helm uninstall notification
helm uninstall rabbitmq
helm uninstall storage

или

make helm_uninstall
```
