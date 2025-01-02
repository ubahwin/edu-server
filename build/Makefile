ifneq (,$(wildcard .env))
    include .env
endif

CONNECTION = $(REMOTE_USERNAME)@$(REMOTE_HOSTNAME)

.PHONY: build-docker
build-docker:
	@docker buildx build --platform linux/amd64 -t edu-server:latest ../. -f Dockerfile

.PHONY: run-docker
run-docker:
	@docker rm -f edu-server-container
	@docker run --name edu-server-container -p 8080:8080 edu-server:latest

.PHONY: docker-image-tar
docker-image-tar: build-docker
	@docker save edu-server:latest > edu-server-image.tar

.PHONY: upload-image-tar
upload-image-tar:
	@echo ">> Uploading image tar"
	@ssh $(CONNECTION) "mkdir -p ~/images"
	@rsync -a --progress edu-server-image.tar $(CONNECTION):~/images
	@ssh $(CONNECTION) "docker rmi -f edu-server"
	@ssh $(CONNECTION) "docker load < ~/images/edu-server-image.tar"

.PHONY: run-remote-docker
run-remote-docker:
	@echo ">> Starting remote docker"
	@ssh $(CONNECTION) "docker rm -f edu-server-container"
	@ssh $(CONNECTION) "docker run --name edu-server-container -p 2346:2345 -p 8081:8080 -p 7070:7070 edu-server:latest"

.PHONY: remote-deploy
remote-deploy: docker-image-tar upload-image-tar upload-compose-config run-remote-docker
