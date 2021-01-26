# Архитектура
![DD](images/SAGA.png)

## Описание сервисов 
| № | Сервис | Описание |Namespace установки|
|---|--------|----------|-------------------|
|1|nginx(ingress)| Точка входа в кластер k8s для потребителей проекта|kube-system|
|2|keycloak| Open ID Conntect Provider авторизации|auth|
|3|oauth2_proxy| Фильтр аутентификации|auth|
|4|profile|Сервис ведения профиля клиента, его инвестеционных потфелей и снимков состояний портфелей по датам (историческая позиция)|default|
|5|PostgreSQL|Коммунальный сервер БД, где для кадого потребителя приходится по дной отдельной БД PostgreSQL|default|
|6|Migration|Набор джобов миграции БД, для каждой БД при раскатке создается и выполняется отдельный JOB миграции|default|
|7|Deals|Сервис ведения сделок(заказов) клиента. Снимок состояния портфеля в сервисе profile ведется по сделкам, которые ведутся в данном сервисе. |default|
|8|notification|Сервис "заглушка" для рассылки уведомлений. В целеом виде предполагается интеграция с шлюзом событий Google|default|
|9|saga-machine|Оркестратор асинхронных бизнес-сценариев (сага)|default|
|12|Kafka|Брокер Kafka|integreation|
|13|ZooKeeper|ZooKeeper|integreation|
## Описание интерфейсов
| № | Сервис | № | Интерфейс | Потребители |Описание интерфейса |API Spec|
|---|--------|---|-----------|-------------|--------------------|--------|
|1|ingress(NGINX)|1.1|arch.homework/auth/|Postman, SPA, Browser|Точка доступа к keycloak для администрирования и получения токенов||
|1|ingress(NGINX)|1.2|arch.homework/profile/api/v1/|Postman, SPA|Точка доступа к сервису управления портфелями profile.default.svc.cluster.local:8080||
|2|Keycloak|2.1|POST /auth/realms/master/protocol/openid-connect/token|Postman, oauth2_proxy|REST API Получение токена пользователя keycloak (implicit flow для posman, code flow  для oauth2_proxy)||
|2|Keycloak|2.2|POST /auth/admin/realms/portfolio/users|Postman|REST API регистрации нового пользователя||
|2|Keycloak|2.3|GET /auth/realms/portfolio|profile<br/>oauth2_proxy</br>deals|Сервис получения конфигурации  OpenID Connet клиента||
|3|oauth2_proxy|3.1|GET /oauth2/auth|ingress(NGINX)|Сервис проверки  JWT  токена или получения JWT токена по сессии oauth2_proxy||
|4|profile|4.1|GET /login|Postman, SPA|Получение данных о пользователе по его токену или сессии oauth2_proxy|[link](https://gitlab.com/portfolio_counselor/profile-src/-/blob/master/profile.yaml)|
|4|profile|4.2-4.5|GET /portfolios<br/>GET /portfolios/{id}<br/>POST /portfolios<br/>PUT /portfolios/{id}|Postman, SPA|CRUD для работы с портфелями|[link](https://gitlab.com/portfolio_counselor/profile-src/-/blob/master/profile.yaml)|
|4|profile|4.6|KAFKA Topic: PROFILE_TOPIC|saga-machine|Топик приема команды на одобрение сделки|Сущность: deal_context<br/>[link](https://gitlab.com/)|
|5|pg|5.1|postgres://keycloak:otus-pg.default.svc.cluster.local:5432/keycloak?sslmode=disable|Keycloak|PostgreSQL БД для Keycloak||
|5|pg|5.2|postgres://profile:otus-pg.default.svc.cluster.local:5432/profile?sslmode=disable|profile|PostgreSQL БД для profile||
|5|pg|5.3|postgres://deals:otus-pg.default.svc.cluster.local:5432/deals?sslmode=disable|deals|PostgreSQL БД для deals||
|5|pg|5.4|postgres://stock:otus-pg.default.svc.cluster.local:5432/stock?sslmode=disable|stock|PostgreSQL БД для stock||
|5|pg|5.5|postgres://notifications:otus-pg.default.svc.cluster.local:5432/notifications?sslmode=disable|notifications|PostgreSQL БД для stock||
|5|pg|5.6|postgres://saga:otus-pg.default.svc.cluster.local:5432/saga?sslmode=disable|saga-machine|PostgreSQL БД для saga-machine||
|6|migration||||Джобы миграции БД||
|7|deals|7.1|POST /deals|SPA, Postman|CRUD сервисы работы со сделкамми|[link](https://gitlab.com/portfolio_counselor/deals-src/-/blob/master/deals.yaml)|
|7|deals|7.2|KAFKA Topic: DEALS_TOPIC|saga-machine|Топик для приема команд по одобрению сделок|Сущность: deal_context<br/>[link](https://gitlab.com/portfolio_counselor/deals-src/-/blob/master/deals.yaml)|
|8|notifications|8.1|KAFKA Topic:  NOTIFICATION_TOPIC|saga-machine|Топик для приема команд по отправке уведомлений|Сущность: notoficationDTO<br/>[link](https://gitlab.com/portfolio_counselor/notification/-/blob/master/notification.yaml)|
|9|notifications|9.1|KAFKA Topic:  SAGA_SERVER_POIC|profile<be/>deals<br/>notification|Топик получения событий от участников саги|Сущность: Event<br/>[link](https://gitlab.com/portfolio_counselor/saga-machine-src/-/blob/master/saga.yaml)|

## Сценарий использования ##
### Статусная машина  MANUAL_CREATE_OREDER (Ручное создание сделки)
![Статусная машина саги](images/MANUAL_CREATE_ORDER.jpg)
**Описание переходов**
|Источник события|Событие от сервиса|Исходное состоение|Следующее состояние|Сервис<br/>получатель комамнды|Команда сервису|Описание|
|----------------|------------------|------------------|-------------------|------------------------------|---------------|--------|
|Deals|REGISTER|-|DRAFT|Profile|APPROVE_OREDER|Сервис сделок регистрирует новую сагу по ручной регистрации сделки|
|Profile|OREDER_APPROVED|DRAFT|APPROVED|Deals|APPROVE_DEAL|Сервис профиля клиента проверяет баланс, одобряет сделку, обновляет состояние портфеля, и возвращает управление оркестратору событием OREDER_APPROVED|
|Profile|OREDER_REJECTED|DRAFT|REJECTED|Deals|REJECT_DEAL|Сервис профиля клиента проверят баланс и отколоняет сделку, и возвращает управление оркестратору событием OREDER_REJECTED|
|Deals|DEAL_UPDATED|PPROVED или REJECTED|NOTIFICATION|notification|NOTIFY|Сервис сделок обновляет информация о сделке, меняет ее статус черновика на Исполнена или Отклонена, и возвращает управление оркестратру командой NOTIFY|

###  Описание шагов сценария ручного создания сделки
![Диаграмма последовательностей](images/deals.png)