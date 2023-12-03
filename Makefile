default: up

env:
	rm -rf ./user_service/.env > /dev/null 2>&1
	cp ./user_service/src/config/parameters/local.env ./user_service/.env

up: down
	docker-compose -f docker-compose.yml up -d

start:
	docker-compose -f docker-compose.yml exec $(service) sh -c \
    	"yarn start:dev"

start-go:
	docker-compose -f docker-compose.yml exec $(service) sh -c \
    	"go run main.go"

down: timeout
	docker ps -a -q | xargs -n 1 -P 8 -I {} docker stop {}
	docker builder prune --all --force
	docker system prune -f

timeout:
	export DOCKER_CLIENT_TIMEOUT=2000
	export COMPOSE_HTTP_TIMEOUT=2000

kill:
	docker-compose -f docker-compose.yml exec -T frontend sh -c \
	"killall node > /dev/null 2>&1"

lint-user_service:
	docker-compose -f docker-compose.yml exec user_service sh -c \
 	 "yarn eslint && yarn prettier && yarn typescript && ANALYZE=true yarn build"

ssh: timeout
	docker-compose -f docker-compose.yml exec $(service) sh

build: down timeout
	docker-compose -f docker-compose.yml build
	docker-compose up -d --remove-orphans
	make package-all
	make down

build-production: down timeout
	docker-compose -f docker-compose.prod.yml build

package-all: 
	make package-user_service

package-user_service:
	docker-compose -f docker-compose.yml exec -T user_service sh -c \
	"yarn install"

owner:
	chown -R $$(whoami) .
	chmod -R 777 .

destroy: down
	docker system prune -a -f

log:
	docker-compose -f docker-compose.yml ps
	sleep 1
	docker-compose -f docker-compose.yml logs -f

bundle-user_service:
	docker-compose -f docker-compose.yml exec -T backend sh -c \
	"yarn build"

# db
migration-generate:
	docker-compose -f docker-compose.yml exec $(service) sh -c \
  "yarn run migration:generate src/migrations/$(name)"

migration-create:
	docker-compose -f docker-compose.yml exec $(service) sh -c \
  "yarn run migration:create src/migrations/$(name)"

migration-run:
	docker-compose -f docker-compose.yml exec -T $(service) sh -c \
	"yarn run migration:run"

migration-revert:
	docker-compose -f docker-compose.yml exec $(service) sh -c \
	"yarn run migration:revert"