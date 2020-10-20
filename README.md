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
helm install otus wwtlf/userlist -f values.yaml  
```

Исходные коды chart лежат тут: https://github.com/WWTLF/userlist

Для тестирования есть коллекция: https://www.getpostman.com/collections/387963098dd12e6e7fd0


# Полезные ссылки
- https://12factor.net/ -  best practices
- https://flagger.app/ - feature toggling
- https://medium.com/hackernoon/wrong-ways-of-defining-service-boundaries-d9e313007bcc - decomposition
- https://skaffold.dev/docs/install/ - CI/CD инструмент
