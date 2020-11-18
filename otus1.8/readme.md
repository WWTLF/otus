# 1.8 Prometheus. Grafana
## Инструментировать сервис из прошлого задания метриками в формате Prometheus с помощью библиотеки для вашего фреймворка и ЯП.

Инжекция метрик:
https://github.com/WWTLF/otus/blob/master/user-list-src/restapi/configure_user_list.go
```
var (
	metrics_rq_latancy = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name: "app_request_latency_seconds",
			Help: "Application Request Latency",
		},
		[]string{"method", "endpoint"},
	)

	metrics_rq_counter = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "app_request_count",
		Help: "Application Request Count",
	}, []string{"method", "endpoint", "http_status"})
)

...

func setupGlobalMiddleware(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		endpoint := r.URL.Path
		if strings.Contains(endpoint, "/user/") {
			endpoint = "/user"
		}
		t0 := time.Now().UnixNano() / (int64(time.Millisecond) / int64(time.Nanosecond))
		w2 := NewLoggingResponseWriter(w)
		handler.ServeHTTP(w2, r)
		total := time.Now().UnixNano()/(int64(time.Millisecond)/int64(time.Nanosecond)) - t0
		metrics_rq_counter.With(prometheus.Labels{"method": r.Method, "endpoint": endpoint, "http_status": fmt.Sprintf("%d", w2.StatusCode)}).Inc()
		metrics_rq_latancy.With(prometheus.Labels{"method": r.Method, "endpoint": endpoint}).Observe(float64(total))
	})
}


```

Демонстрация, что метрики отдаются:

![counter](https://github.com/WWTLF/otus/blob/master/otus1.8/counter.png)

![latency](https://github.com/WWTLF/otus/blob/master/otus1.8/Hist.png)


## Сделать дашборд в Графане, в котором были бы метрики с разбивкой по API методам:
  - [x] 1. Latency (response time) с квантилями по 0.5, 0.95, 0.99, max
  - [x] 2. RPS
  - [x] 3. Error Rate - количество 500ых ответов
  
 ![main](https://github.com/WWTLF/otus/blob/master/otus1.8/main_dash_board.png)

## Добавить в дашборд графики с метрикам в целом по сервису, взятые с nginx-ingress-controller:
  - [x] 1. Latency (response time) с квантилями по 0.5, 0.95, 0.99, max
  - [x] 2. RPS
  - [x] 3. Error Rate - количество 500ых ответов

Скриншоты с алертингом ниже.
  
## Настроить алертинг в графане на Error Rate и Latency.

- На выходе должно быть:
 - [x] 0) скриншоты дашборды с графиками в момент стресс-тестирования сервиса. Например, после 5-10 минут нагрузки.
 
![nginx](https://github.com/WWTLF/otus/blob/master/otus1.8/NGINX_dash_board.png)
 
 - [x] 1) json-дашборды:
   - Для сервиса https://github.com/WWTLF/otus/blob/master/otus1.8/grafana-configmap.yaml 
   - Для INGRESS https://github.com/WWTLF/otus/blob/master/otus1.8/grafana-nginx-configmap.yaml 

## Задание со звездочкой (+5 баллов)
### Используя существующие системные метрики из кубернетеса, добавить на дашборд графики с метриками:
  - [x] 1. Потребление подами приложения памяти
![cpu](https://github.com/WWTLF/otus/blob/master/otus1.8/cpu.png)
  - [x] 2. Потребление подами приолжения CPU
![mem](https://github.com/WWTLF/otus/blob/master/otus1.8/mem.png)

Релизовано с помощью встроенного дашборда: kubernetes-compute-resources-pod

  

### Инструментировать базу данных с помощью экспортера для prometheus для этой БД.
- [x] Добавить в общий дашборд графики с метриками работы БД.
  - Ативируем сбор метрик в настройках values чарта: https://github.com/WWTLF/userlist/blob/master/values.yaml
```
pg:
  metrics:
    enabled: true
    serviceMonitor:
      enabled: true
```
 - Настраиваем еще один дашборд в графане: https://github.com/WWTLF/otus/blob/master/otus1.8/grafana-pg-configmap.yaml
 
 ![pg](https://github.com/WWTLF/otus/blob/master/otus1.8/pg.png)
