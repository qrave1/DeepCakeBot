# DeepCake Bot

Универсальный помощник для Telegram с функцией утренней рассылки прогноза погоды.


### Запуск локально (для разработки)

Запустите инфраструктуру:

```bash
docker-compose up -d
```

Отредактируйте `.env` и заполните необходимые значения:

- `TELEGRAM_BOT_TOKEN` - токен бота от [@BotFather](https://t.me/BotFather)
- `OPENWEATHER_API_KEY` - API ключ от [OpenWeather](https://openweathermap.org/api)
- `DATABASE_URL` - строка подключения к PostgreSQL

Запустите бота:

```bash
go run main.go
```

## Переменные окружения

| Переменная | Описание | Значение по умолчанию |
|-----------|----------|----------------------|
| `TELEGRAM_BOT_TOKEN` | Токен Telegram бота | - (обязательно) |
| `OPENWEATHER_API_KEY` | API ключ OpenWeather | - (обязательно) |
| `DATABASE_URL` | URL подключения к PostgreSQL | - (обязательно) |
| `TIMEZONE` | Часовой пояс для планировщика | Europe/Moscow |
| `WEATHER_SCHEDULE_HOUR` | Час отправки прогноза (0-23) | 7 |
| `CITY` | Город для прогноза погоды | Moscow |
| `COUNTRY_CODE` | Код страны (ISO 3166) | RU |

