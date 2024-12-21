# Run postgresql container
docker run -d \
--name accrual_db \
-e POSTGRES_USER=accrual \
-e POSTGRES_PASSWORD='!qaz2wsx' \
-e POSTGRES_DB=accrual \
-v "$(pwd)/pgdata/accrual:/var/lib/postgresql/data" \
-p 5433:5432 \
postgres:latest