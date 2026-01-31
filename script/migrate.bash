#! /bin/sh
# Script to control postgres container migrations from backend container
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

DB_HOST="simple_go_backend-dev-db-1"
DB_PASSWORD=$(cat "$DB_PASSWORD_FILE")
SSLMODE=${SSLMODE:-disable}

DB_URL="postgres://$DB_USER:$DB_PASSWORD@$DB_HOST:$DB_PORT/$DB_NAME?sslmode=$SSLMODE"

if [ "$#" = 0 ]; then
  /go/bin/migrate -source file:///go/db/migrations/ -database "${DB_URL}" up 
elif [ "$1" = "down" ]; then
  /go/bin/migrate -source file:///go/db/migrations/ -database "${DB_URL}" down
else
  /go/bin/migrate -source file:///go/db/migrations/ -database "${DB_URL}" up "$1"
fi;
