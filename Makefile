.PHONY: run imports fmt db db-rm view

# General
APP_ENV=develop

# Src vars
SRC=main.go pdf/*.go config/*.go db/*.go parsers/*.go
MAIN_SRC=main.go

# Db vars
DB_DOCKER_NAME=postgres-${APP_ENV}
DB_PORT=5432
DB_USER=${APP_ENV}
DB_PASSWORD=${APP_ENV}
DB_NAME=crime-map-${APP_ENV}

# Runs the server
run:
	go run ${MAIN_SRC}

# Adds all required go imports
imports:
	goimports -l -w ${SRC}

# Formats go source
fmt:
	gofmt -w ${SRC}

# View all source files
view:
	${EDITOR} $(SRC)

# Runs the postgres db
db: 
	docker run \
		--name ${DB_DOCKER_NAME} \
		--net=host \
		-p ${DB_PORT}:${DB_PORT} \
		-e POSTGRES_USER=${DB_USER} \
		-e POSTGRES_PASSWORD=${DB_PASSWORD} \
		-e POSTGRES_DB=${DB_NAME} \
		postgres \
	|| \
	docker start -ai $(shell docker ps -qaf name=${DB_DOCKER_NAME})

# Deletes the postgres db
db-rm:
	docker rm ${DB_DOCKER_NAME}
