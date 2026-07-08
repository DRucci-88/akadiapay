package domain

type AppConfigProvider interface {
	APP_PORT() int
	APP_ENV() string
	DB_DSN() string
	JWT_SECRET() string
}
