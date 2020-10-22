## Задание 1.6. Основы работы с Kubernetes (часть 3)

Для установки:
```
helm repo add wwtlf https://wwtlf.github.io/userlist
helm repo update
helm install otus wwtlf/userlist -f values.yaml  
```

Для тестирования в данной папке есть коллекция POSTMAN:
```
newman run otusapp-borisershov.postman_collection.json
```

Исходные коды chart лежат тут: https://github.com/WWTLF/userlist
