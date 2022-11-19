# Metrics collector

Учебный проект агент-серверного сборщика метрик (WIP)

На данный момент агент собирает не имеющие смысла на практике собственные runtime метрики, которые периодически пушит на сервер.
Для создания "практического" сборщика метрик необходимо реализовать интерфейс
```go
type Poller interface {
	Poll(*metric.Metrics)
}
```
и заинжектить его в ``pollworker`` (инжект в cmd/agent/main.go)

Поддержка нескольких агентов подразумевает отправку метрик с идентификатором агентской машины в имени метрики.
