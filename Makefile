PHONY: build run build-run

# Docker Composeでコンテナを起動するターゲット
run:
	docker-compose up --build