## Выполнение домашнего задания

Сервис написан на golang

### Команды для сборки и запуска сервиса

- `make docker-build` - сборка docker image с сервисом
- `make docker-push` - отправка docker image в docker registry
- `make compose-up` - запуск docker-compose с сервисом
- `make compose-down` - остановка docker-compose с сервисом

### Проверка

```
# запуск контейнера
docker run --name health -d -p 8000:8000 astrviktor/health:1.0.0

# запрос
curl --request GET 'http://127.0.0.1:8000/health/'

# ответ
{"status":"OK"}

# остановка контейнера
docker stop health

# удаление контейнера
docker rm health
```

