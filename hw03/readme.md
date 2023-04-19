## Выполнение домашнего задания

### Установка minikube

- https://www.linuxtechi.com/how-to-install-minikube-on-ubuntu/

### Установка helm

- https://helm.sh/docs/intro/install/

```
curl https://baltocdn.com/helm/signing.asc | gpg --dearmor | sudo tee /usr/share/keyrings/helm.gpg > /dev/null
sudo apt-get install apt-transport-https --yes
echo "deb [arch=$(dpkg --print-architecture) signed-by=/usr/share/keyrings/helm.gpg] https://baltocdn.com/helm/stable/debian/ all main" | sudo tee /etc/apt/sources.list.d/helm-stable-debian.list
sudo apt-get update
sudo apt-get install helm

helm version
```


### Старт minikube с ingress

```
minikube start --addons=ingress --cpus=2 --cni=flannel --install-addons=true --kubernetes-version=stable --memory=8g
```

### Простая установка postgres в minikube через helm

- https://habr.com/ru/companies/domclick/articles/649167/
- https://fusionauth.io/docs/v1/tech/installation-guide/kubernetes/setup/minikube

```
helm repo list

kubectl create namespace storage
helm install pg-minikube --set auth.postgresPassword=password bitnami/postgresql --namespace storage

пароль сохранился в секретах
kubectl get secret --namespace storage pg-minikube-postgresql -o jsonpath="{.data.postgres-password}" | base64 -d
```

### Установка postgres в minikube через helm с PersistentVolume

- https://phoenixnap.com/kb/postgresql-kubernetes

```
kubectl apply -f ./deployments/kubernetes/storage/storage-namespace.yml
kubectl apply -f ./deployments/kubernetes/storage/storage-pv.yml
kubectl apply -f ./deployments/kubernetes/storage/storage-pvc.yml

kubectl get pvc -n storage

helm install storage bitnami/postgresql \
  --set persistence.existingClaim=postgresql-pv-claim \
  --set volumePermissions.enabled=true \
  --set auth.postgresPassword=password \
  --set auth.database=users \
  --namespace storage

kubectl port-forward --namespace storage svc/storage-postgresql 5432:5432

Первоначальные миграции не требуются
```

### Создание secret для postgres
```
user: postgres
password: password

echo -n 'postgres' | base64
echo -n 'password' | base64

cG9zdGdyZXM=
cGFzc3dvcmQ=

apiVersion: v1
kind: Secret
metadata:
  name: crudservice-secret
  namespace: crud  
type: Opaque
data:
  user: cG9zdGdyZXM=
  password: cGFzc3dvcmQ=
  
```

### Запуск сервиса crudservice в Kubernetes

```
kubectl apply -f ./deployments/kubernetes/crud/crud-namespace.yml
kubectl apply -f ./deployments/kubernetes/crud/crud-secret.yml
kubectl apply -f ./deployments/kubernetes/crud/crud-configmap.yml
kubectl apply -f ./deployments/kubernetes/crud/crud-deployment.yml
kubectl apply -f ./deployments/kubernetes/crud/crud-service.yml
kubectl apply -f ./deployments/kubernetes/crud/crud-ingress.yml
```

### Проверка Ingress, Service, Deployment

```
kubectl get ingress -n crud

NAME             CLASS   HOSTS           ADDRESS        PORTS   AGE
health-ingress   nginx   arch.homework   192.168.49.2   80      49s

kubectl get all -n crud

NAME                              READY   STATUS    RESTARTS   AGE
pod/crudservice-8bf4b5c8b-htlrh   1/1     Running   0          36s
pod/crudservice-8bf4b5c8b-ttp6l   1/1     Running   0          35s

NAME                  TYPE       CLUSTER-IP       EXTERNAL-IP   PORT(S)          AGE
service/crudservice   NodePort   10.101.239.233   <none>        8000:30001/TCP   30s

NAME                          READY   UP-TO-DATE   AVAILABLE   AGE
deployment.apps/crudservice   2/2     2            2           36s

NAME                                    DESIRED   CURRENT   READY   AGE
replicaset.apps/crudservice-8bf4b5c8b   2         2         2       36s
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

---
# запрос
curl --request POST 'http://arch.homework/user' --data '{"username": "johndoe589", "firstName": "John", "lastName": "Doe", "email": "bestjohn@doe.com", "phone": "+71002003040"}'

# ответ
{"id":1}

---
# запрос
curl --request GET 'http://arch.homework/user/1'

# ответ
{"id":1,"username":"johndoe589","firstName":"John","lastName":"Doe","email":"bestjohn@doe.com","phone":"+71002003040"}

---
# запрос
curl --request PUT 'http://arch.homework/user/1' --data '{"username": "johndoe666"}'

# ответ
{"id":1}

---
# запрос
curl --request GET 'http://arch.homework/user/1'

# ответ
{"id":1,"username":"johndoe666","firstName":"John","lastName":"Doe","email":"bestjohn@doe.com","phone":"+71002003040"}

---
# запрос
curl --request DELETE 'http://arch.homework/user/1'

# ответ
{"id":1}

---
# запрос
curl --request GET 'http://arch.homework/user/1'

# ответ
record not found
```

### Удаление minikube
```
minikube stop
minikube delete
```
