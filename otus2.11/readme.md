### Шаг 1. Устновка  главного приложения из HELM

### Шаг 2. Устновка  keycloak

```
kubectl create namespace auth
helm repo add codecentric https://codecentric.github.io/helm-charts
helm repo update
helm install keycloak codecentric/keycloak -f keycloak-values.yaml -n auth
```

*Если имя релиза главного приложения отличается от otus, то необходимо указать правильно ссылку на БД в файле keycloak-valyues.yaml
```
- name: DB_ADDR
    value: otus-pg.default.svc.cluster.local
```
Сервисы keycloak и главного приложения делят один сервер БД, но разные инстансы БД. 

### Шаг 3. Устновка  oauth2_proxy

```
kubectl apply -f oauth2-deployment.yaml
```