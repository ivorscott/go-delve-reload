#!make

include .env

GOROOT=$(shell go env GOROOT)
NETWORKS="$(shell docker network ls)"
VOLUMES="$(shell docker volume ls)"
SCHEMA_DIR=api/internal/schema
SEED_DIR=$(SCHEMA_DIR)/seeds
MIGRATION_DIR=$(SCHEMA_DIR)/migrations
MIGRATIONS_VOLUME= $(PWD)/$(MIGRATION_DIR):/migrations
USER_PASS_HOST=$(POSTGRES_USER):$(POSTGRES_PASSWORD)@$(POSTGRES_HOST)
URL=postgres://$(USER_PASS_HOST):5432/$(POSTGRES_DB)?sslmode=disable
SUCCESS=[ done "\xE2\x9C\x94" ]

# This ALLOWS the following usage -> "make migration <name>", "make seed <name>", "make insert <name>"
# Normally, this is not the case in Makefiles, usually it's -> "make migrations name=<name>" etc. http://bit.ly/37TR1r2
ifeq ($(firstword $(MAKECMDGOALS)),$(filter $(firstword $(MAKECMDGOALS)),migration seed insert))
  name := $(wordlist 2,$(words $(MAKECMDGOALS)),$(MAKECMDGOALS))
  $(eval $(name):;@:)
endif
# This ALLOWS the following usage -> "make up <number>", "make down <number>"
ifeq ($(firstword $(MAKECMDGOALS)),$(filter $(firstword $(MAKECMDGOALS)),up down force))
  num := $(wordlist 2,$(words $(MAKECMDGOALS)),$(MAKECMDGOALS))
  $(eval $(num):;@:)
# When we migrate down without a number the number defaults to 1.
# In other words, if we do "make down" rather than "make down <number>".
# Note: migrating up without a number "make up" has no default on purpose to migrate to the latest migration.
# Therefore, you must specifically provide a number to prevent this -- "make up <number>"
  ifndef num
    ifeq ($(firstword $(MAKECMDGOALS)),$(filter $(firstword $(MAKECMDGOALS)),down))
      num := 1
    endif
  endif
endif

# default arguments
user ?= root
service ?= api

cert:
	@echo [ generating self-signed certificates ... ]
	@echo
	@mkdir -p ./api/tls \
	&& go run $(GOROOT)/src/crypto/tls/generate_cert.go --rsa-bits=2048 --host=localhost \
	&& mv *.pem ./api/tls \
	&& echo \
	&& echo file: api/tls/cert.pem \
	&& echo file: api/tls/key.pem \
	&& echo \
	&& echo ============================================= \
	&& echo  Getting Started Docs: http://bit.ly/38hLFWG \
	&& echo ============================================= \
	&& echo \
	&& echo $(SUCCESS)

postgres-network:
ifeq (,$(findstring $(POSTGRES_NET),$(NETWORKS)))
	@echo [ creating postgres network... ]
	@docker network create $(POSTGRES_NET)
	@echo $(SUCCESS)
endif
	
api: postgres-network
	@echo [ starting api ... ]
	@docker-compose up api 

db: postgres-network
	@echo [ running postgres in the background ... ]
	@docker-compose up -d db
	@docker-compose ps   

client:
	@echo [ starting client ... ]
	@docker-compose up client

clean:
	@echo [ teardown all containers ... ]
	docker-compose down
	@echo $(SUCCESS)

audit:
	@echo [ audit client container ... ]
	@make exec service="client" cmd="npm audit"
	@echo $(SUCCESS)

audit-fix:
	@echo [ audit fix client container ... ]
	@make exec service="client" cmd="npm audit fix"
	@echo $(SUCCESS)

tidy: 
	@echo [ cleaning up unused $(service) dependencies ... ]
	@make exec service="api" cmd="go mod tidy"

exec:
	@echo [ executing $(cmd) in $(service) ]
	docker-compose exec -u $(user) $(service) $(cmd)
	@echo $(SUCCESS)

test-client:
	@echo [ running client tests ... ]
	@cd client; npx jest --notify
	@cd ..

test-client-watch:
	@echo [ running client tests ... ]
	@cd client; npx jest --watchAll
	@cd ..

test-api:
	@echo [ running api tests ... ]
	@cd api; go test -v ./...
	@cd ..

debug-api:
	@echo [ debugging api ... ]
	@docker-compose up debug-api

debug-db:
	@echo [ debugging db ... ]
	@echo
	@# advanced command line interface for postgres
	@# includes auto-completion and syntax highlighting. https://www.pgcli.com/
	@docker run -it --rm --net $(POSTGRES_NET) dencold/pgcli $(URL)

rm:
	@echo [ removing all containers ... ]
	docker rm -f `docker ps -aq`

rmi:
	@echo [ removing all images ... ]
	docker rmi -f `docker images -a -q`

migration:
    ifndef name
		$(error migration name is missing -> make migration <name>)
    endif

	@echo [ generating migration files ... ]
	@docker run --volume $(MIGRATIONS_VOLUME) --network $(POSTGRES_NET) migrate/migrate \
	create \
	-ext sql \
	-dir /migrations \
	-seq $(name) \
	&& echo \
	&& echo located at $(MIGRATION_DIR) \
	&& echo \
	&& echo migrations \
	&& echo \
	&& ls api/internal/schema/migrations \
	&& echo \
	&& echo $(SUCCESS)

version: 
	@echo [ printing migration version ... ]
	@echo
	@docker run --volume $(MIGRATIONS_VOLUME) --network $(POSTGRES_NET) migrate/migrate \
	-path /migrations \
	-database $(URL) version \
	&& echo \
	&& echo  $(SUCCESS)
	
up:
	@echo [ migrating up ... ]
	@echo
	@docker run --volume $(MIGRATIONS_VOLUME) --network $(POSTGRES_NET) migrate/migrate \
	-path /migrations \
	-verbose \
	-database $(URL) up $(num) \
	&& echo \
	&& echo $(SUCCESS)

down:
	@echo [ migrating down ... ]
	@echo
	@docker run --volume $(MIGRATIONS_VOLUME) --network $(POSTGRES_NET) migrate/migrate \
	-path /migrations  \
	-verbose \
	-database $(URL) down $(num) \
	&& echo \
	&& echo $(SUCCESS)

force: 
	@echo [ forcing version ... ]
	@echo
# A migration script can fail because of invalid syntax in sql files. http://bit.ly/2HQHx5s
# To fix this, force a previous version to be the new current one
# 1) "make force <version>"
# 2) fix the syntax issue
# 3) then run "make up" again
	@docker run --volume $(MIGRATIONS_VOLUME) --network $(POSTGRES_NET) migrate/migrate \
	-path /migrations \
	-verbose \
	-database $(URL) force $(num) \
	&& echo \
	&& echo $(SUCCESS)

seed:
    ifndef name
		$(error seed name is missing -> make insert <name>)
    endif
	
	@echo [ generating seed file ... ]
	@echo
	@echo file: $(SEED_DIR)/$(name).sql
	@echo
	@mkdir -p $(PWD)/$(SEED_DIR)
	@touch $(PWD)/$(SEED_DIR)/$(name).sql \
	&& echo $(SUCCESS)

insert:
    ifndef name
		$(error seed filename is missing -> make insert <filename>)
    endif

	@echo [ inserting $(name) seed data ... ]
	@echo
	@docker cp $(PWD)/$(SEED_DIR)/$(name).sql $(shell docker-compose ps -q db):/seed/$(name).sql \
	&& docker exec -u root db psql $(POSTGRES_DB) $(POSTGRES_USER) -f /seed/$(name).sql \
	&& echo \
	&& echo $(SUCCESS)

.PHONY: all
.PHONY: api
.PHONY: cert
.PHONY: client
.PHONY: exec
.PHONY: db
.PHONY: debug-api
.PHONY: debug-db
.PHONY: rm
.PHONY: rmi
.PHONY: down
.PHONY: dump
.PHONY: force
.PHONY: insert
.PHONY: migration
.PHONY: postgres-network
.PHONY: teardown
.PHONY: test-client
.PHONY: test-api
.PHONY: tidy
.PHONY: seed
.PHONY: up
.PHONY: version
