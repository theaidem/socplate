include .env

define check_arg
	@[ ${1} ] || ( echo ">> ${2} is not set, use: ${2}=value"; exit 1 )
endef

SERVER_BIN='server/${APP_NAME}'
DATABASE_URL='postgres://localhost:${DATABASE_PORT}/${DATABASE_DB}?sslmode=disable&user=${DATABASE_USER}&password=${DATABASE_PASSWORD}'
MIGRATIONS_PATH=$${PWD}/migrations
DATABASE_DRIVER="host=database user=$${POSTGRES_USER} dbname=$${POSTGRES_DB} password=$${POSTGRES_PASSWORD} sslmode=disable"

test:	
	@echo ${APP_NAME} ${arg}

appname:	
	@echo Application server named: ${APP_NAME} 

build.client: 
	@cd client && yarn build
	@echo build.client: OK!

build.server.local: clean.server
	@cd server && go build -o server
	@echo build.server.local: OK!

build.server.linux: clean.server
	@cd server && GOOS=linux GOARCH=amd64 go build -o server
	@echo build.server.linux: OK!

start.dev: appname build.server.linux
	@docker-compose -f docker-compose.yaml -f docker-compose.dev.yaml up --build

start.prod: appname build.server.linux build.client
	# TODO prod env
	@docker-compose -f docker-compose.yaml -f docker-compose.prod.yaml up -d --build

clean.server:
	@rm -rf ${SERVER_BIN}
	@echo clean.server: OK!

migartion.create:
	$(call check_arg, ${name}, name)
	@docker-compose exec migrations goose -v -dir /migrations create ${name} sql
	@echo 'Migration ${name} created'

migration.status:
	@docker-compose exec migrations ash -c \
	'goose -v -dir /migrations postgres ${DATABASE_DRIVER} status'
	
migration.up:
	@docker-compose exec migrations ash -c \
	'goose -v -dir /migrations postgres ${DATABASE_DRIVER} up'

migration.down:
	@docker-compose exec migrations ash -c \
	'goose -v -dir /migrations postgres ${DATABASE_DRIVER} down'