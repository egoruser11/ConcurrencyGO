# Concurrency GO

<p align="center">
  Набор учебных задач по конкурентности на Go с аккуратной структурой, понятными названиями и примерами базовых паттернов работы с goroutine, channel, select, context и worker pool.
</p>

<p align="center">
  <img alt="Go" src="https://img.shields.io/badge/Go-1.18-00ADD8?style=for-the-badge&logo=go&logoColor=white">
  <img alt="Patterns" src="https://img.shields.io/badge/Patterns-10-success?style=for-the-badge">
  <img alt="Status" src="https://img.shields.io/badge/Status-Learning%20Project-orange?style=for-the-badge">
</p>

## О проекте

Этот репозиторий собран как практическая коллекция задач по конкурентному программированию на Go.
Каждая задача вынесена в отдельную папку в корне проекта и запускается независимо через свой `main.go`.

Основная цель проекта:

- потренировать работу с `goroutine` и `channel`
- понять, как строятся конвейеры обработки данных
- разобрать паттерны `fan-in`, `fan-out`, `worker pool`
- научиться использовать `context` для корректной остановки горутин
- увидеть примеры `select`, retry, rate limit и параллельных продюсеров

## Что внутри

| Папка | Паттерн | Описание |
| --- | --- | --- |
| `basic_pipeline_pattern` | Базовый pipeline | Генерация чисел, обработка и чтение результата в простом конвейере. |
| `concurrent_producers_pattern` | Несколько продюсеров | Несколько горутин одновременно пишут в один общий канал. |
| `context_cancellation_pattern` | Graceful shutdown | Остановка фоновой работы через `context.CancelFunc`. |
| `context_filter_pipeline_pattern` | Pipeline + context | Генерация, преобразование, фильтрация и сбор данных в одной цепочке. |
| `fan_in_fan_out_pattern` | Fan-out + fan-in | Распределение задач по нескольким обработчикам и объединение результата. |
| `fan_in_square_pattern` | Fan-in | Слияние нескольких входных потоков в один выходной канал. |
| `miner_pool_pattern` | Worker pool + context | Несколько шахтёров параллельно добывают ресурс и отдают его в общий канал. |
| `priority_select_pattern` | Приоритетный `select` | Чтение из нескольких каналов с разным приоритетом. |
| `processing_pipeline_pattern` | Многоэтапный pipeline | Обработка структуры данных через несколько стадий. |
| `rate_limited_worker_pattern` | Worker pool + retry | Ограничение частоты обработки, повтор запросов и сбор результатов. |

## Структура проекта

```text
.
├── README.md
├── go.mod
├── basic_pipeline_pattern/
├── concurrent_producers_pattern/
├── context_cancellation_pattern/
├── context_filter_pipeline_pattern/
├── fan_in_fan_out_pattern/
├── fan_in_square_pattern/
├── miner_pool_pattern/
├── priority_select_pattern/
├── processing_pipeline_pattern/
└── rate_limited_worker_pattern/
```

## Быстрый старт

### Требования

- Go 1.18+

### Запуск любого примера

```bash
go run ./basic_pipeline_pattern
```

Или, например:

```bash
go run ./miner_pool_pattern
go run ./fan_in_fan_out_pattern
go run ./rate_limited_worker_pattern
```

## Идея оформления

В проекте сделан упор не только на рабочий код, но и на читаемость:

- каждая задача находится в отдельной директории
- в начале каждого файла есть комментарий с описанием паттерна
- функции и переменные названы по смыслу
- примеры разделены по логике, чтобы их было проще читать и запускать по отдельности

## Для чего подойдёт этот репозиторий

- как учебный конспект по конкурентности в Go
- как база для повторения перед собеседованиями
- как набор простых шаблонов для собственных экспериментов
- как понятная коллекция примеров для дальнейшего расширения

## Дальше можно развить

- добавить тесты для отдельных паттернов
- собрать общий индекс всех задач по темам
- добавить схемы потоков данных для каждого примера
- оформить отдельные примеры с `mutex`, `errgroup` и bounded worker pool
