# Переменные
APP_NAME=inventory-service
DOCKER_COMPOSE=docker-compose.yml

.PHONY: help build up down restart logs ps test fmt clean

# По умолчанию выводит список команд
help:
	@echo "Использование: make [команда]"
	@echo "Команды:"
	@echo "  build   - Сборка Docker образов"
	@echo "  up      - Запуск всех контейнеров в фоне"
	@echo "  down    - Остановка и удаление контейнеров"
	@echo "  restart - Перезапуск сервисов"
	@echo "  logs    - Просмотр логов приложения"
	@echo "  ps      - Статус контейнеров"
	@echo "  test    - Запуск тестов"
	@echo "  fmt     - Форматирование кода (go fmt)"
	@echo "  clean   - Очистка бинарников и кэша Docker"

build:
	docker-compose build

up:
	docker-compose up -d

down:
	docker-compose down

restart: down up

logs:
	docker-compose logs -f app

ps:
	docker-compose ps

test:
	go test -v -race ./...

fmt:
	go fmt ./...
	go mod tidy

clean:
	rm -f server
	docker system prune -f