#/bin/bash
docker run -d -p 6379:6379 redis

docker run -d -p 3306:3306 -e MYSQL_ROOT_PASSWORD=password mysql:8.0.30

docker run -d -p 5432:5432 -e POSTGRES_PASSWORD=password postgres:9.6
