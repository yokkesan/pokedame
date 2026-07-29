package database

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net"
	"net/url"
	"os"
	"strconv"
	"sync"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
)

const (
	defaultMaxOpenConnections = 10
	defaultMaxIdleConnections = 5
	defaultConnectionLifetime = 30 * time.Minute
	connectionTimeout         = 5 * time.Second
)

var (
	db   *sql.DB
	dbMu sync.RWMutex
)

func Initialize() error {
	connection, err := open()
	if err != nil {
		return err
	}

	dbMu.Lock()
	defer dbMu.Unlock()

	if db != nil {
		_ = connection.Close()

		return errors.New("データベースは既に初期化されています")
	}

	db = connection

	return nil
}

func Get() (*sql.DB, error) {
	dbMu.RLock()
	defer dbMu.RUnlock()

	if db == nil {
		return nil, errors.New("データベースが初期化されていません")
	}

	return db, nil
}

func Close() error {
	dbMu.Lock()
	defer dbMu.Unlock()

	if db == nil {
		return nil
	}

	if err := db.Close(); err != nil {
		return fmt.Errorf("データベース接続を終了できませんでした: %w", err)
	}

	db = nil

	return nil
}

func open() (*sql.DB, error) {
	dsn, err := buildDSN()
	if err != nil {
		return nil, err
	}

	connection, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, fmt.Errorf("データベースを初期化できませんでした: %w", err)
	}

	connection.SetMaxOpenConns(defaultMaxOpenConnections)
	connection.SetMaxIdleConns(defaultMaxIdleConnections)
	connection.SetConnMaxLifetime(defaultConnectionLifetime)

	ctx, cancel := context.WithTimeout(context.Background(), connectionTimeout)
	defer cancel()

	if err := connection.PingContext(ctx); err != nil {
		_ = connection.Close()

		return nil, fmt.Errorf("データベースへ接続できませんでした: %w", err)
	}

	return connection, nil
}

func buildDSN() (string, error) {
	host := os.Getenv("DB_HOST")
	portValue := os.Getenv("DB_PORT")
	name := os.Getenv("DB_NAME")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")

	if host == "" || portValue == "" || name == "" || user == "" || password == "" {
		return "", errors.New("データベース接続用の環境変数が不足しています")
	}

	port, err := strconv.Atoi(portValue)
	if err != nil || port < 1 || port > 65535 {
		return "", errors.New("DB_PORTが不正です")
	}

	connectionURL := &url.URL{
		Scheme: "postgres",
		User:   url.UserPassword(user, password),
		Host:   net.JoinHostPort(host, portValue),
		Path:   name,
	}

	query := connectionURL.Query()
	query.Set("sslmode", resolveSSLMode())
	connectionURL.RawQuery = query.Encode()

	return connectionURL.String(), nil
}

func resolveSSLMode() string {
	if os.Getenv("APP_ENV") == "local" {
		return "disable"
	}

	return "require"
}
