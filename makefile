#!make

NETWORK="$(shell docker network ls | grep traefik-public)"

network:
ifneq (,$(findstring traefik-public,$(NETWORK)))
    # Network already exists
else
    # Network doesn't exist
	@echo
	@echo [ creating network... ]
	docker network create traefik-public
	@echo [ done ]

endif

api: network
	@echo
	@echo [ starting api... ]
	docker-compose up traefik api

api-d:
	@echo
	@echo [ teardown api... ]
	docker-compose down
	@echo [ done ]

debug-api:
	@echo
	@echo [ starting debug-api... ]
	docker-compose up debug-api

run:
	@echo
	@echo [ executing $(cmd) in new api container... ]
	docker-compose run -u root --rm api $(cmd)
	@echo [ done ]

exec:
	@echo
	@echo [ executing $(cmd) in running api container... ]
	docker-compose exec -u root api $(cmd)
	@echo [ done ]

test:
	@make exec cmd="go test ./..."

.PHONY:	api
