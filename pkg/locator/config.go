package locator

type Config struct {
	Port      int
	Namespace string
}

func (c *Config) Defaults() {
	if c.Port == 0 {
		c.Port = 8080
	}
	if c.Namespace == "" {
		c.Namespace = "default"
	}
}
