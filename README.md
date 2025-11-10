# vk-response-preparer

## Что делает

1. Загружает конфигурацию из окружения, валидирует обязательные параметры Kafka и создает Zap-логгер.
2. Создает `sarama`-producer с `SCRAM-SHA512`, идемпотентной отправкой и компрессией Snappy для темы, куда пишутся готовые VKMessages (`TOPIC_NAME_VK_REQUEST_MESSAGE`).
3. Формирует `KafkaConsumer` в группе `GROUP_ID_VK_RESPONSE_PREPARER`, подписываясь на тему `TOPIC_NAME_VK_RESPONSE_PREPARER` и снабжает обработчик `processor.TgMessagePreparer`.
4. Обработчик десериализует `contract.NormalizedResponse`, берет `chat_id` и `text`, формирует `contract.SendMessageRequest` и пересылает его в продьюсер `TOPIC_NAME_VK_REQUEST_MESSAGE`.
5. Все ошибки логируются, consumer держит пул сигналов для корректного завершения на `SIGINT/SIGTERM`.

## Запуск

1. Подготовьте `.env` (см. следующую секцию) с нужными переменными.
2. Соберите и запустите в контексте Go:
   ```bash
   go run ./cmd/vk-response-preparer
   ```
3. Либо постройте Docker-образ и запустите, передав файл `.env`:
   ```bash
   docker build -t vk-response-preparer .
   docker run --rm --env-file .env vk-response-preparer
   ```
4. Для локального отладки удобно экспортировать переменные:
   ```bash
   set -a && source .env && set +a && go run ./cmd/vk-response-preparer
   ```

## Переменные окружения

### Kafka (все обязательны, кроме SASL, если не используется)

- `KAFKA_BOOTSTRAP_SERVERS_VALUE` — `host:port[,host:port]` списка брокеров.
- `KAFKA_TOPIC_NAME_VK_REQUEST_MESSAGE` — топик, куда отправляются подготовленные `SendMessageRequest` для отправки в VK.
- `KAFKA_TOPIC_NAME_VK_RESPONSE_PREPARER` — входной топик, к которому подписан consumer.
- `KAFKA_GROUP_ID_VK_RESPONSE_PREPARER` — идентификатор consumer group.
- `KAFKA_CLIENT_ID_VK_RESPONSE_PREPARER` — идентификатор клиента (producer & consumer).
- `KAFKA_SASL_USERNAME` и `KAFKA_SASL_PASSWORD` — передайте, если кластер требует SASL/PLAIN.

## Примечания

- JSON-сообщения, которые потребляются, должны соответствовать `internal/contract.NormalizedResponse`.
- `processor.TgMessagePreparer` не обрабатывает пустой текст: если `text` пустой, он все равно отправит `message:""` в итоговый топик.
- Логика граничит на том, что входящий `NormalizedResponse` уже содержит `chat_id`/`text`, нужные `source`/`context` используются на более высоких уровнях сервиса.
