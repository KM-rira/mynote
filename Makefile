.PHONY: restart
restart:
	docker compose down
	docker compose up --build -d

.PHONY: start
start:
	docker compose up --build -d

.PHONY: down
down:
	docker compose down


