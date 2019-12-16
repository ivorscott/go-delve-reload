#!make

rm:
	@echo "\n[ stopping and removing docker containers ]"
	@docker stop `docker ps -a -q`
	@docker rm `docker ps -a -q`
	@echo "\n"

rmi:
	@echo "\n[ removing docker images ]"
	docker rmi -f `docker images -a -q`
	@echo "\n"

api:
	@echo "\n[ startup api]"
	@docker-compose up --build api
	@echo "\n"

debug-api:
	@echo "\n[ debug api ]"
	@docker-compose up --build debug-api
	@echo "\n"

.PHONY:	api
