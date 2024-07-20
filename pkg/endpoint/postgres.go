package endpoint

import "fmt"

const (
	SslModeDisabled = "disable"

	userDefault     = "postgres"
	passwordDefault = "postgres"
	hostDefault     = "localhost"
	portDefault     = 5432
	dbDefault       = "postgres"
	sslModeDefault  = SslModeDisabled
)

type options struct {
	user string
	pass string
	host string
	port int
	db   string
	ssl  string
}

type PostgresOption func(o *options)

func BuildPostgresConnectionString(opts ...PostgresOption) string {
	options := &options{
		user: userDefault,
		pass: passwordDefault,
		host: hostDefault,
		port: portDefault,
		db:   dbDefault,
		ssl:  sslModeDefault,
	}

	for _, opt := range opts {
		opt(options)
	}

	return fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=%s",
		options.user, options.pass, options.host, options.port, options.db, options.ssl,
	)
}

func WithPostgresUser(user string) PostgresOption {
	return func(o *options) {
		o.user = user
	}
}

func WithPostgresPassword(password string) PostgresOption {
	return func(o *options) {
		o.pass = password
	}
}

func WithPostgresHost(host string) PostgresOption {
	return func(o *options) {
		o.host = host
	}
}

func WithPostgresPort(port int) PostgresOption {
	return func(o *options) {
		o.port = port
	}
}

func WithPostgresDb(db string) PostgresOption {
	return func(o *options) {
		o.db = db
	}
}
