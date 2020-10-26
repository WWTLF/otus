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

### Краткое описание

- [x] PostgreSQL подключен как зависимость
- [x] Геренрируется 2 Job'а:
  - create-user-job -  создает нового пользователя Postgres с ограниченными правами
  - userlist-migrate Производит миграцию данных с помощью инстурмента https://github.com/golang-migrate/migrate, исходный код контейнера миграции тут: https://github.com/WWTLF/otus/tree/master/userlist-migration-src
- [x] Главный сервис userlist https://github.com/WWTLF/otus/tree/master/user-list-src
- [x] Пароли лежат в secrets
- [x] Точка монтирования PV на хосте определяется настройками БД по умолчанию: /tmp/hostpath-provisioner
  - Есть возможность указать свой PVC в файле volumes в разделе pg
- [x] Настроена liveness проба
- [x] Настроена readiness проба, которая проверяет доступность БД
- [x] Реализован gracefull shutdown, который безопасно завершает соединение с БД
