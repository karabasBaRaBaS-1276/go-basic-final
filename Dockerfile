# 1. Шаг сборки
# Базовый образ golang версии 1.24.2
FROM golang:1.24.2 AS builder
# Рабочая дирректория относительно которой будет все операции будут выполняться
WORKDIR /opt/go-app/src
# Копируем все файлы нужные для сборки (меняются редко)
COPY go.mod go.sum ./
# Скачать зависимости в соответствии с go.mod
RUN go mod download
# Копируем все файлы проекта (меняются часто) с учетом исключений из .dockerignore
COPY . ./
# Cобрать приложение go для запуска
# Удалить все исходники относительно рабочей дирректории, кроме статики в папке web
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /opt/go-app/go-basic-final \
    && find . -maxdepth 1 ! -name 'web' ! -name '.' -exec rm -rf {} +
# Создать системную (--system) группу gouser
# Создать системного (--system) пользователя gouser в группе (--ingroup gouser),
#         без возможности входа через оболочку (--shell /bin/false), 
#         без пароля (--disabled-password) и домашней дирректории (--no-create-home)
# Меняем владельца рекурсивно с отчетом
# Указанные действия делаем тут, чтобы в scratch образе (фаза 2) не заниматься установкой 
#         недостающих команд для работы с пользователями и группами
RUN addgroup --system gouser \
    && adduser --system --ingroup gouser --shell /bin/false --no-create-home --disabled-password gouser \
    && chown -v -R gouser:gouser /opt/go-app

# 2. Шаг релиза образа
FROM scratch
WORKDIR /opt/go-app/src
# Указываем здесь ARG, ENV и EXPOSE
ARG TODO_HOST="0.0.0.0" \
    TODO_PORT=7540 \
    TODO_DBDRIVER="sqlite" \
    TODO_DBFILE="scheduler.db"
ENV TODO_HOST=${TODO_HOST} \
    TODO_PORT=${TODO_PORT} \
    TODO_DBDRIVER=${TODO_DBDRIVER} \
    TODO_DBFILE=${TODO_DBFILE}

EXPOSE ${TODO_PORT}
# Копируем из builder приложение с установленными правами на него
COPY --from=builder /opt/go-app /opt/go-app
# Копируем из builder пользователей и группы
COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /etc/group /etc/group

# Приложение будет запускаться под пользователем gouser
USER gouser
# Запустить приложение
CMD ["/opt/go-app/go-basic-final"]