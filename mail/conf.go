package mail

import (
	"github.com/spf13/viper"
	"sync"
)

type emailConfig struct {
	host     string
	port     int
	name     string
	from     string
	password string
}

var emailConf *emailConfig
var once sync.Once

func MustStartup() *emailConfig {
	once.Do(func() {
		emailConf = newEmailConfig()
	})
	return emailConf
}

func newEmailConfig() *emailConfig {
	host := viper.GetString("email.host")
	port := viper.GetInt("email.port")
	name := viper.GetString("email.name")
	from := viper.GetString("email.from")
	password := viper.GetString("email.password")

	config := &emailConfig{
		host:     host,
		port:     port,
		name:     name,
		from:     from,
		password: password,
	}
	return config
}
