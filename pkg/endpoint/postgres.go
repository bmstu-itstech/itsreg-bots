package endpoint

import "fmt"

const (
	SslModeDisabled = "disable"
)

func PostgresConnectionString(
	host string,
	port string,
	user string,
	pass string,
	name string,
	mode string,
) string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		host, port, user, pass, name, mode)
}
