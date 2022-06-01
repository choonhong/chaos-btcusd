.PHONY: all
all: help

SHELL:=/bin/bash

FILE=./.env
ifneq ("$(wildcard $(FILE))","")
	include $(FILE)
	export $(shell sed 's/=.*//' $(FILE))
endif

.PHONY: help
help:
	###############################################################################
	#
	# available commands:
	#
	# up            - run the application and database in docker containers
	# run           - run the app locally
	# down          - run docker-compose down (shutdown all the containers)
	#
	###############################################################################

# Up all the services
.PHONY: up
up: 
	docker-compose -f build/app/docker-compose.yaml up --build -d

.PHONY: down
down:
	docker-compose -f build/app/docker-compose.yaml down

.PHONY: run
run:
	go mod tidy && go run cmd/run/main.go
