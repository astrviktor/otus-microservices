## Выполнение домашнего задания

Разработан отдельный сервис заказов `orderservice`

Метод POST /order сделан идемпотетным

Идемпотетность работает через header X-Request-Id

Работает следующим образом:
1) Если запрос на создание приходит без X-Request-Id, сервис генерирует X-Request-Id = UUID к заказу и сохраняет данные в базе
2) Если запрос на создание приходит с X-Request-Id, сервис ищет заказ в базе по X-Request-Id
   - если находит, то возвращает имеющийся в базе заказ
   - если не находит, то создает в базе новый заказ


### Установка сервиса orderservice и postgres (namespace default)

```
helm install storage ./deployments/kubernetes/helm-charts/postgres/ \
    --set auth.postgresPassword=password \
    --set auth.database=orders

helm install orderservice ./deployments/kubernetes/helm-charts/orderservice/

или

make helm_install
```

### Проверка тестового сценария Postman (newman)
```
newman run ./postman/hw07-orders.postman_collection.json --verbose

или

make newman_run
```


### Удаление сервиса orderservice и postgres
```
helm uninstall orderservice
helm uninstall storage

или

make helm_uninstall
```
