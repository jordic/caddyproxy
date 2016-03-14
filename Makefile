
RELEASE=1

build:
	CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main github.com/jordic/caddyproxy

docker_build:

	docker rmi jordic/caddy-gen; docker build -t jordic/caddy-gen:$(RELEASE) .

push:
	docker push jordic/caddy-gen
