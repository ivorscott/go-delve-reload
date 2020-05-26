#!make
# USE THIS FILE FOR DEPLOYMENT ONLY

include .env

build: 
	@echo "\n[ build api production image ]"

	docker build --target prod \
	--build-arg version=v1 \
	--build-arg backend=${REACT_APP_BACKEND} \
	--tag devpies/gdr-client ./client

	docker build --target prod \
	--tag devpies/gdr-api ./api

login: 
	@echo "\n[ log into private registry ]"
	cat ./secrets/registry_pass | docker login --username `cat ./secrets/registry_user` --password-stdin

publish:
	@echo "\n[ publish production grade images ]"
	docker push devpies/gdr-api
	docker push devpies/gdr-client

deploy:
	@echo "\n[ startup production stack ]"
	@cat ./startup
	@docker stack deploy -c docker-stack.yml --with-registry-auth gdr

metrics: 
	@echo "\n[ enable docker engine metrics ]"
	./init/enable-monitoring.sh

secrets: 
	@echo "\n[ create swarm secrets ]"
	./init/create-secrets.sh

servers:
	@echo "\n[ create servers ]"
	./init/create-servers.sh

servers-d:
	@echo "\n[ teardown swarm ]"
	./init/destroy-servers.sh

swarm:
	@echo "\n[ create swarm with all managers ]"
	./init/create-swarm.sh

.PHONY: build 
.PHONY: login
.PHONY: publish
.PHONY: deploy
.PHONY: metrics
.PHONY: secrets
.PHONY: servers
.PHONY: servers-d
.PHONY: swarm
