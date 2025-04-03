run:
	docker compose --env-file .env up --build

logs-all:
	docker compose logs -f

logs-raft:
	docker compose logs order_service_2 order_service_1 order_service_3 -f

logs_queue:
	docker compose logs order_queue -f

kill:
	docker compose restart order_service_${ID}
