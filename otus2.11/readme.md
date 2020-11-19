- [Otus 2.11](#otus-211)
  * [Архитектура](#-----------)
  * [Инструкция по установке](#-----------------------)
    + [Шаг 1. Устновка  главного приложения из HELM](#----1----------------------------------helm)
    + [Шаг 2. Устновка  keycloak](#----2-----------keycloak)
    + [Шаг 3. Устновка  oauth2_proxy](#----3-----------oauth2-proxy)
  * [Описание тестовых сценариев](#---------------------------)

<small><i><a href='http://ecotrust-canada.github.io/markdown-toc/'>Table of contents generated with markdown-toc</a></i></small>


# Otus 2.11 

## Архитектура

## Инструкция по установке

### Шаг 1. Устновка  главного приложения из HELM

helm install nginx ingress-nginx/ingress-nginx -f nginx-ingress.yaml

### Шаг 2. Устновка  keycloak

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


**Шаг 2.1 Настройка  CoreDNS:**

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

### Шаг 3. Устновка  oauth2_proxy

```
kubectl apply -f oauth2-deployment.yaml
```

## Описание тестовых сценариев
