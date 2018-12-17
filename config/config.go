package config

//Config config
type Config struct {
	FindLimit       int
	MongoDbURL      string
	MongoDbDatabase string
}

var config *Config = nil

// Get get
func Get() *Config {
	if config == nil {
		config = &Config{FindLimit: 1000, MongoDbURL: "127.0.0.1", MongoDbDatabase: "crud"}
	}
	return config
}
