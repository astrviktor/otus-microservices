## Выполнение домашнего задания

1-2 пункт выполнены в ДЗ1

### Установка minikube

- https://www.linuxtechi.com/how-to-install-minikube-on-ubuntu/

### Старт minikube с ingress

```
minikube start --addons=ingress --cpus=2 --cni=flannel --install-addons=true --kubernetes-version=stable --memory=8g
```

### Запуск в Kubernetes

```
kubectl apply -f .
```

### Проверка Ingress

```
kubectl get ingress -n health

NAME             CLASS   HOSTS           ADDRESS        PORTS   AGE
health-ingress   nginx   arch.homework   192.168.49.2   80      49s
```

### Добавление в /etc/hosts

```
sudo nano /etc/hosts

192.168.49.2    arch.homework
```

### Проверка
```
# запрос
curl --request GET 'http://arch.homework/health/'

# ответ
{"status":"OK"}

# запрос
curl --request GET 'http://arch.homework/otusapp/astrviktor/health/'

# ответ
{"status":"OK"}

```

### Удаление minikube
```
minikube stop
minikube delete
```
