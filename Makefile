run:
	docker compose down --volumes
	docker compose up

build:
	docker compose down --volumes
	docker compose build

test:
	make clean
	docker compose -f docker-compose.test.yml up --build --remove-orphans backend_test_viewer
	docker compose -f docker-compose.test.yml up --build --remove-orphans frontend_test_viewer

clean:
	docker compose -f docker-compose.test.yml down
	docker compose -f docker-compose.test.yml down --volumes
	docker compose down
	docker compose down --volumes