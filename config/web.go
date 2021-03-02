package config

type Web struct {
	port string
}

func (w Web) GetPort() string {
	return w.port
}

