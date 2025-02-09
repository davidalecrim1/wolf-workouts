package config

type Config struct {
	jwtSecret []byte
}

func NewConfig(jwtSecret string) *Config {
	return &Config{
		jwtSecret: []byte(jwtSecret),
	}
}

func (c *Config) GetJWTSecret() []byte {
	return c.jwtSecret
}
