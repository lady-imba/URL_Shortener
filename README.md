# URL Shortener

Сервис для сокращения URL с авторизацией и редиректами.

## Как запустить

### 1. Клонируйте репозиторий
git clone https://github.com/ВАШ_НИК/URL_Shortener.git
cd URL_Shortener

### 2. Настройте конфигурацию
cp config/local.yaml.example config/local.yaml
Отредактируйте config/local.yaml под свои нужды

### 3. Запустите сервер
go run cmd/url-shortener/main.go

## Немного про API
Требуют Basic Auth
POST /url - создать короткую ссылку
DELETE /url/{alias} - удалить ссылку

Публичные
GET /{alias} - редирект на оригинальный URL

### Технологии
Go 
Chi Router
SQLite
Basic Auth
