NAME=simple_brocker

SOURCE=cmd/app/main.go

PACKAGE=\
	github.com/labstack/echo/v4\
	github.com/ilyakaznacheev/cleanenv\

all:
	go run $(SOURCE)
