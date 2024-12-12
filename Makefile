include .env

BINARY_NAME=server

build:
	go build -o ${BINARY_NAME} cmd/main.go


migration-up:
	cd ./migrations
	goose sqlite3 ../storage.db up 
	cd ..

migration-down:
	echo $(PWD)	
	cd ./migrations
	goose sqlite3 ../storage.db down
	cd ..

seed:
	export CONFIG_PATH=${CONFIG_PATH}
	go run seedData.go

run:
	export CONFIG_PATH=${CONFIG_PATH}
	./${BINARY_NAME}
