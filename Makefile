IMAGE_NAME = go-proxy
CONTAINER_NAME = go-proxy

build:
	docker build -t $(IMAGE_NAME) .

run:
	docker run -p 8080:8080 --rm -v $(PWD):/go/src/go-proxy-server/ --env-file .env -it --name $(CONTAINER_NAME) $(IMAGE_NAME)

inspect:
	docker inspect $(CONTAINER_NAME)

shell:
	docker exec -it $(CONTAINER_NAME) /bin/sh

stop:
	docker stop $(CONTAINER_NAME)
