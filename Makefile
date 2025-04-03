run:
	docker compose --env-file .env up --build

logs:
	docker compose logs -f

