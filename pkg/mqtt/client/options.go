package client

type Option func(*client)

func UID(username string) Option {
	return func(s *client) {
		s.uid = username
	}
}

func Password(password string) Option {
	return func(s *client) {
		s.password = password
	}
}

func Host(host string) Option {
	return func(s *client) {
		s.host = host
	}
}

func Port(port string) Option {
	return func(s *client) {
		s.port = port
	}
}

func AddLogger(base Logger) Option {
	return func(c *client) {
		c.logger = NewLogger(base)
	}
}
