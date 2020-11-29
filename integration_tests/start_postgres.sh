export DB_USERNAME=postgres
export DB_PASSWORD=test
export DB_SCHEMA=features

docker run --rm \
--name test_postgres \
-e POSTGRES_USER=$DB_USERNAME \
-e POSTGRES_PASSWORD=$DB_PASSWORD \
-e POSTGRES_DB=$DB_SCHEMA \
-p 54300:5432 \
postgres:11


