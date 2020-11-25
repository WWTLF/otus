- [Otus 2.11](#otus-211)
  * [Архитектура](#Архитектура)
  * [Инструкция по установке](#Инструкция-по-установке)
    * [Шаг 1. Устновка главного приложения из HELM](#Шаг-1--Устновка-главного-приложения-из-HELM)
    * [Шаг 2. Устновка keycloak](#Шаг-2--Устновка-keycloak)
    * [Шаг 3. Устновка oauth2_proxy](#Шаг-3--Устновка-oauth2_proxy)
  * [Описание тестовых сценариев](#Описание-тестовых-сценариев)


# Otus 2.11 

## Архитектура

![DD](images/portfolio.png)

| № | Серивис | № | Интерфейс | Описание интерфейса |
|---|---------|---|-----------|---------------------|
|1|ingress(NGINX)|1.1|arch.homework/auth/|Точка доступа к keycloak для администрирования и получения токенов|
|1|ingress(NGINX)|1.2|arch.homework/profile/api/v1/|Точка доступа к сервису упревления портфелями|
|2|Keycloak|2.1|POST /auth/realms/master/protocol/openid-connect/token|REST API Получение токена пользователя keycloak (implicit flow)|
|2|Keycloak|2.2|POST /auth/admin/realms/portfolio/users|REST API регистрации нового пользователя|
|2|Keycloak|2.3|GET /auth/realms/portfolio|Сервис получения конфигурации  OpenID Connet клиента|
|3|oauth2_proxy|3.1|GET /oauth2/auth|Сервис проверки  JWT  токена или получения JWT токена по сессии oauth2_proxy|

## Инструкция по установке

### Пререквизиты

1. Должен быть установлен ingress

```
minikube addons enable ingress
```

2. Необходимо склонировать данный репозиторий и открыть папку otus2.11 со всеми необходимыми файлами для настройки сервисов

```
git clone git@github.com:WWTLF/otus.git
cd otus2.11
```

### Шаг 1. Устновка главного приложения из HELM

```
helm repo add wwtlf2 https://wwtlf.github.io/portfolio-chart/
helm repo update
helm install otus wwtlf2/portfolio-chart
```

### Шаг 2. Устновка keycloak

```
kubectl create namespace auth
helm repo add codecentric https://codecentric.github.io/helm-charts
helm repo update
helm install keycloak codecentric/keycloak -f keycloak-values.yaml -n auth
```
*HELM главного приложения несет с собой БД для keycloak с настроенной конфигурацией*

*Если имя релиза главного приложения отличается от otus, то необходимо указать правильно ссылку на БД в файле keycloak-valyues.yaml*
```
- name: DB_ADDR
    value: otus-pg.default.svc.cluster.local
```
Сервисы keycloak и главного приложения делят один сервер БД, но разные инстансы БД. 


**Шаг 2.1 Настройка CoreDNS:**

Стандарт OIDC требует, чтобы хост получения токена совпадал с хостом получения ключей для проверки подписи токена. Поэтому доступ к keycloak, как через Ingress, так и через Service(ClusterIP) должен идти по хосту arch.homework. В k8s нет стандартного ресурса для добавления host alias к сервису (можно только к IP, но это не гибко), но можно добавить DNS правило для CoreDNS Controller. 

```
kubectl apply -f coredns-configmap.yaml
# Перезапускаем coredns
kubectl get pods -n kube-system | grep coredns
coredns-f9fd979d6-j5z4c        1/1     Running   0          45m
# Указываем правильный POD coredsn-...
kubectl delete pod coredns-f9fd979d6-j5z4c -n kube-system
```

*(FYI) Содержимое файла coredns-configmap.yaml:*
```
apiVersion: v1
data:
  Corefile: |
    .:53 {
        errors
        health {
           lameduck 5s
        }
        ready
        kubernetes cluster.local in-addr.arpa ip6.arpa {
           pods insecure
           fallthrough in-addr.arpa ip6.arpa
           ttl 30
        }
        rewrite name arch.homework keycloak-http.auth.svc.cluster.local //ВАЖНО! Правильно указать внешний хост и FQDN сервиса Keycloak
        prometheus :9153
        forward . /etc/resolv.conf {
           max_concurrent 1000
        }
        cache 30
        loop
        reload
        loadbalance
    }
kind: ConfigMap
metadata:
  name: coredns
  namespace: kube-system
  
```

### Шаг 3. Устновка oauth2_proxy

```
kubectl apply -f oauth2-deployment.yaml
```

## Описание тестовых сценариев
