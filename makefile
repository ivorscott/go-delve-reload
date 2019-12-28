#!make

include .env

SEED_DIR=api/internal/schema/seeds
MIGRATIONS_VOLUME= $(PWD)/api/internal/schema/migrations:/migrations
USER_PASS_HOST=$(POSTGRES_USER):$(POSTGRES_PASSWORD)@$(POSTGRES_HOST)
URL=postgres://$(USER_PASS_HOST):5432/$(POSTGRES_DB)?sslmode=disable

# This ALLOWS the following usage -> "make migration <name>", "make seed <name>", "make insert <name>"
# Normally, it would be -> "make migrations name=<name>" etc. http://bit.ly/37TR1r2
ifeq ($(firstword $(MAKECMDGOALS)),$(filter $(firstword $(MAKECMDGOALS)),migration seed insert))
  name := $(wordlist 2,$(words $(MAKECMDGOALS)),$(MAKECMDGOALS))
  $(eval $(name):;@:)
endif
# This ALLOWS the following usage -> "make up <number>", "make down <number>"
ifeq ($(firstword $(MAKECMDGOALS)),$(filter $(firstword $(MAKECMDGOALS)),up down force))
  num := $(wordlist 2,$(words $(MAKECMDGOALS)),$(MAKECMDGOALS))
  $(eval $(num):;@:)
# When we migrate down without a number the number defaults to 1.
  ifndef num
    ifeq ($(firstword $(MAKECMDGOALS)),$(filter $(firstword $(MAKECMDGOALS)),down))
      num := 1
    endif
  endif
endif

migration:
	@echo [ generating migration files ... ]
    ifndef name
		$(error migration name is missing -> make migration <name>)
    endif
	@docker run --volume $(MIGRATIONS_VOLUME) --network $(POSTGRES_NET) migrate/migrate \
	create -ext sql -dir /migrations -seq $(name)

version: 
	@echo [ printing migration version ... ]
	@docker run --volume $(MIGRATIONS_VOLUME) --network $(POSTGRES_NET) migrate/migrate \
	-path /migrations \
	-database $(URL) version
up:
	@echo [ migrating up ... ]
	@docker run --volume $(MIGRATIONS_VOLUME) --network $(POSTGRES_NET) migrate/migrate \
	-path /migrations -verbose -database $(URL) up $(num)

down:
	@echo [ migrating down ... ]
	@docker run --volume $(MIGRATIONS_VOLUME) --network $(POSTGRES_NET) migrate/migrate \
	-path /migrations -verbose -database $(URL) down $(num)

force: 
	@echo [ forcing version ... ]
	# A migration script can fail because of invalid syntax in sql files. http://bit.ly/2HQHx5s
	@docker run --volume $(MIGRATIONS_VOLUME) --network $(POSTGRES_NET) migrate/migrate \
	-path /migrations -verbose -database $(URL) force $(num)

seed:
	@echo [ generating seed file ... ]
    ifndef name
		$(error seed name is missing -> make insert <name>)
    endif
	@mkdir -p $(PWD)/$(SEED_DIR);touch $(PWD)/$(SEED_DIR)/$(name).sql 

insert:
	@echo [ inserting $(name) seed data ... ]
    ifndef name
		$(error seed filename is missing -> make insert <filename>)
    endif
	@docker cp $(PWD)/$(SEED_DIR)/$(name).sql $(shell docker-compose ps -q db):/seed/$(name).sql \
	&& docker exec -u root db psql $(POSTGRES_DB) $(POSTGRES_USER) -f /seed/$(name).sql 

.PHONY: all
.PHONY: cert
.PHONY: down
.PHONY: force
.PHONY: insert
.PHONY: migration
.PHONY: version
.PHONY: seed
.PHONY: up