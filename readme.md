docker run --name web-service-exercise-mysql -e MYSQL_ROOT_PASSWORD=q1w2e3r4 -p 3306:3306 -d public.ecr.aws/docker/library/mysql:8.4
docker run --name web-service-exercise-postgres -e POSTGRES_PASSWORD=q1w2e3r4 -p 5432:5432 -d public.ecr.aws/docker/library/postgres:17.4

curl http://localhost:8080/albums \
--include \
--header "Content-Type: application/json" \
--request "POST" \
--data '{"title": "The Modern Sound of Betty Carter","artist": "Betty Carter","price": 49.99}'

GOOS=linux GOARCH=arm64 go build -tags lambda.norpc -o bootstrap main.go 