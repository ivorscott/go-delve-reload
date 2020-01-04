#!make

NETWORKS="$(shell docker network ls)"
POSTGRES_DB="$(shell cat ./secrets/postgres_db)"
POSTGRES_HOST="$(shell cat ./secrets/postgres_host)"
POSTGRES_USER="$(shell cat ./secrets/postgres_user)"
POSTGRES_PASSWORD="$(shell cat ./secrets/postgres_passwd)"
SUCCESS=[ done "\xE2\x9C\x94" ]

# default arguments
user ?= root
service ?= api

all:
	@echo [ starting client '&' api... ]
	docker-compose up traefik client api db pgadmin

traefik-network:
ifeq (,$(findstring traefik-public,$(NETWORKS)))
	@echo [ creating traefik network... ]
	docker network create traefik-public
	@echo $(SUCCESS)
endif

postgres-network:
ifeq (,$(findstring postgres,$(NETWORKS)))
	@echo [ creating postgres network... ]
	docker network create postgres
	@echo $(SUCCESS)
endif

api: traefik-network postgres-network
	@echo [ starting api... ]
	docker-compose up traefik api db pgadmin

down:
	@echo [ teardown all containers... ]
	docker-compose down
	@echo $(SUCCESS)

install:
	@echo [ installing $(service) dependencies... ]
	@make exec cmd="$(cmd)"

tidy: 
	@echo [ cleaning up unused $(service) dependencies... ]
	@make exec cmd="go mod tidy"

exec:
	@echo [ executing $(cmd) ]
	docker-compose exec -u $(user) $(service) $(cmd)
	@echo $(SUCCESS)

test:
	@echo [ running tests... ]
	@make exec cmd="go test -v ./..."

debug-api:
	@echo [ debugging api... ]
	docker-compose up traefik debug-api db pgadmin

debug-db:
	@echo [ debugging postgres database... ]
	@# basic command line interface for postgres 
	@# make exec user="$(POSTGRES_USER)" service="$(POSTGRES_HOST)" cmd="bash -c 'psql --dbname $(POSTGRES_DB)'"

	@# advanced command line interface for postgres
	@# includes auto-completion and syntax highlighting. https://www.pgcli.com/
	@docker run -it --rm --net postgres dencold/pgcli postgresql://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@$(POSTGRES_HOST):5432/$(POSTGRES_DB)

dump:
	@echo [ dumping postgres backup for $(POSTGRES_DB)... ]
	@docker exec -it $(POSTGRES_HOST) pg_dump --username $(POSTGRES_USER) $(POSTGRES_DB) > ./api/scripts/backup.sql
	@echo $(SUCCESS)

.PHONY:	api
