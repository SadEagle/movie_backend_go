package main

import (
	"context"
	"fmt"
	"os"
	"strconv"

	"github.com/jackc/pgx/v5"
)

func insertAdminUser(conn *pgx.Conn) error {
	query := `
        INSERT INTO user_data(name, login, encoded_password, is_admin) 
        SELECT 'Sr. Admin', 'admin', $1::bytea, true 
        WHERE NOT EXISTS(
            SELECT NULL
            FROM user_data
            WHERE login = 'admin'
        )`

	// Hex-encoded password as a byte array
	encodedPassword := []byte{0x19, 0x36, 0x14, 0xb0, 0xf1, 0xfd, 0x19, 0x44, 0xc2, 0xd3, 0xc3, 0x85, 0x36, 0x68, 0x60, 0x39, 0x80, 0x4d, 0x33, 0xd1, 0x9a, 0x01, 0x38, 0x72, 0x55, 0x2b, 0xa7, 0x71, 0xc9, 0x50, 0x30, 0x7d}
	ctx := context.TODO()
	_, err := conn.Exec(ctx, query, encodedPassword)
	return err
}

func main() {

	port, err := strconv.Atoi(os.Getenv("DB_PORT"))
	if err != nil {
		panic("cant detect port")
	}
	text, err := os.ReadFile(os.Getenv("DB_PASSWORD_FILE"))
	if err != nil {
		panic("cant detect password")
	}

	Host := os.Getenv("DB_HOST")
	Port := port
	Database := os.Getenv("DB_NAME")
	User := os.Getenv("DB_USER")
	Password := string(text)
	SSLMode := "disable"
	connUrl := fmt.Sprintf("host=%s port=%d dbname=%s user=%s password=%s sslmode=%s", Host, Port, Database, User, Password, SSLMode)

	ctx := context.TODO()
	conn, err := pgx.Connect(ctx, connUrl)
	if err != nil {
		panic("connection error")
	}
	insertAdminUser(conn)
}
