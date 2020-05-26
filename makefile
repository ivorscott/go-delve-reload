#!make
# DON'T EXTEND THE SCOPE OF THIS FILE FOR DEVELOPMENT CONCERNS
# USE THIS FILE FOR DEPLOYMENT ONLY

include .env

build: 
	@echo "\n[ building production api images ]"

	docker build --target prod \
	--build-arg version=v1 \
	--build-arg backend=${REACT_APP_BACKEND} \
	--tag devpies/gdr-client ./client

	docker build --target prod \
	--tag devpies/gdr-api ./api

login: 
	@echo "\n[ logging into private registry ]"
	cat ./secrets/registry_pass | docker login --username `cat ./secrets/registry_user` --password-stdin

publish:
	@echo "\n[ publishing production grade images ]"
	docker push devpies/gdr-api
	docker push devpies/gdr-client

deploy:
	@echo "\n[ deploying production stack ]"
	@cat ./startup
	@docker stack deploy -c docker-stack.yml --with-registry-auth gdr

metrics: 
	@echo "\n[ enabling docker engine metrics ]"
	./init/enable-monitoring.sh

secrets: 
	@echo "\n[ creating swarm secrets ]"
	./init/create-secrets.sh

server:
	@echo "\n[ creating server ]"
	./init/create-server.sh

server-d:
	@echo "\n[ destroying server ]"
	./init/destroy-server.sh

swarm:
	@echo "\n[ create single node swarm ]"
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
