# otus
Для домашних заданй
## Задание 1.4. Основы работы с Kubernetes (часть 2)
Папка с манифестами: https://github.com/WWTLF/otus/tree/master/otus1.4

Для проверки:
```
curl http://arch.homework/otusapp/borisershov/health
```

## Задание 1.5. Основы работы с Kubernetes (часть 3)

### Способ 1. БД из HELM, Сервисы из kubectl apply -f:

Для установки БД:
```
helm repo add bitnami https://charts.bitnami.com/bitnami
helm install user-list-db bitnami/postgresql 
```
Далее применение манифестов: 
```
kubectl apply -f deployment.yaml
kubectl apply -f service.yaml
kubectl apply -f ingress.yaml
```
TODO: Добавить манифесты секретов

### Способ 2: Установка с помощью HELM
