export DB_USERNAME=mongo
export DB_PASSWORD=test
export DB_SCHEMA=features


docker run -d \
-p 27017:27017 \
-e MONGO_INITDB_ROOT_USERNAME=$DB_USERNAME \
-e MONGO_INITDB_ROOT_PASSWORD=$DB_PASSWORD \
-e MONGO_INITDB_DATABASE=$DB_SCHEMA \
mongo:latest

# --network some-network --name some-mongo