## Функции, которые реализованы:

### Функция copyMap

    func copyMap(metricList *map[string]*uint64) map[string]uint64

+ копирует мапу и возращает ее с измененным типом значения по ключу
+ тип значения изменяется с *uint64 --> uint64

### Функция updateAndCount

     func updateAndCount(m *rt.MemStats, metricList *map[string]*uint64, pollCounter *counter) 

+ обновляет метрики
+ считает количество изменившихся

### Функция createMetricList

    func createMetricList(pollCounter *counter) (metricList map[string]*uint64)

+ создает объект типа 'runtime.MemStates'
+ созает мапу с метриками

### Функция sendData 

    func sendData(metricsList *map[string]*uint64)

+ конструирует url
+ отправляет метрики на адресс 127.0.0.1:8080

### Функция MyMonitor

    func MyMpnitor()

+ запускает бесконечный цикл обновления и отправки метрик
+ раз в 2 секунды обновляет
+ раз в 10 секунд отправляет

выполнено при помощи двух тикеров из модуля 'time'