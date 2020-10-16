# otus
Для домашних заданй
## Задание 1.4. Основы работы с Kubernetes (часть 2)
Папка с манифестами: https://github.com/WWTLF/otus/tree/master/otus1.4

Для проверки:
```
curl http://arch.homework/otusapp/borisershov/health
```

## Задание 1.6. Основы работы с Kubernetes (часть 3)

Для установки:
```
helm repo add wwtlf https://wwtlf.github.io/userlist
helm repo update
helm install otus userlist -f values.yaml  
```

Исходные коды лежат тут: https://github.com/WWTLF/userlist


# Полезные ссылки
- https://12factor.net/ -  best practices
- https://flagger.app/ - feature toggling
- https://stackoverflow.com/questions/50218376/managing-db-migrations-on-kubernetes-cluster - database migrations
- https://medium.com/hackernoon/wrong-ways-of-defining-service-boundaries-d9e313007bcc - decomposition
