run:
	docker compose --env-file .env up --build

logs_all:
	docker compose logs -f

logs_raft:
	docker compose logs order_service_2 order_service_1 order_service_3 -f

logs_queue:
	docker compose logs order_queue -f

kill:
	docker compose restart order_service_${ID}

run_simple:
	docker compose --env-file .env -f docker docker-compose.simple.yaml up --build