#! /bin/sh
cd /migrate
go mod init w
go get github.com/jackc/pgx/v5

go run init_admin.go
