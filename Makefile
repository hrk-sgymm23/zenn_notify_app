PHONY: build run build-run

init:
	cd src
	go mod init zenn_notify_app

run:
	docker-compose up --build