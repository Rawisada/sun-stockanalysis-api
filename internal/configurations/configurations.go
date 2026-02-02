package configurations

import (
	"strings"
	"sync"
	"time"

	"github.com/spf13/viper"
	"github.com/go-playground/validator/v10"
	
)

type (
	Config struct {
		Server   	*Server   	`mapstructure:"server" validate:"required"`
		// OAuth2   	*OAuth2   	`mapstructure:"oauth2" validate:"required"`
		State   	*State   	`mapstructure:"state" validate:"required"`
		Database 	*Database 	`mapstructure:"database" validate:"required"`
		Finnhub		*Finnhub	`mapstructure:"finnhub" validate:"required"`
	}

	Server struct {
		Port           	int      		`mapstructure:"port" validate:"required"`
		AllowedOrigins 	[]string 		`mapstructure:"allowOrigins" validate:"required"`
		BodyLimit		string			`mapstructure:"bodyLimit" validate:"required"`
		TimeOut        	time.Duration  	`mapstructure:"timeout" validate:"required"`
		
	}

	// OAuth2 struct {
	// 	ClientID     	string 		`mapstructure:"clientId" validate:"required"`
	// 	ClientSecret 	string 		`mapstructure:"clientSecret" validate:"required"`
	// 	RedirectURL  	string 		`mapstructure:"redirectUrl" validate:"required"`
	// 	EndPoints		endpoint 	`mapstructure:"endpoints" validate:"required"`
	// 	Scopes			[]string	`mapstructure:"scopes" validate:"required"`
	// 	UserInfoUrl		string 		`mapstructure:"userInfoUrl" validate:"required"`
	// 	revokeUrl		string 		`mapstructure:"revokeUrl" validate:"required"`
	// }

	// endpoint struct {
	// 	AuthUrl     	string 		`mapstructure:"authUrl" validate:"required"`
	// 	TokenUrl 		string 		`mapstructure:"tokenUrl" validate:"required"`
	// 	DeviceAuthUrl  	string 		`mapstructure:"deviceAuthUrl" validate:"required"`
	// }

	State    struct{
		Secret     		string 				`mapstructure:"secret" validate:"required"`
		ExpiredsAt 		time.Duration 		`mapstructure:"expiredsAt" validate:"required"`
		Issuer  		string 				`mapstructure:"issuer" validate:"required"`
	}

	Database struct {
		Host     		string 		`mapstructure:"host" validate:"required"`
		Port     		int    		`mapstructure:"port" validate:"required"`
		User     		string 		`mapstructure:"user" validate:"required"`
		Password 		string 		`mapstructure:"password" validate:"required"`
		DBname   		string 		`mapstructure:"dbname" validate:"required"`
		SSLmode   		string 		`mapstructure:"sslmode" validate:"required"`
		Schema   		string 		`mapstructure:"schema" validate:"required"`
	}

	Finnhub struct {
		Token string `mapstructure:"token" validate:"required"`
	}
)

var (
	once 				sync.Once  //Singleton 
	configInstance 		*Config
)

func ConfigGetting() *Config {
	once.Do(func() {
		viper.SetConfigName("config")
		viper.SetConfigType("yaml")
		viper.AddConfigPath("./config")
		viper.AddConfigPath(".")
		viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
		viper.AutomaticEnv()

		if err := viper.ReadInConfig(); err != nil {
			panic(err)
		}

		var cfg Config
		if err := viper.Unmarshal(&cfg); err != nil {
			panic(err)
		}

		if err := validator.New().Struct(&cfg); err != nil {
			panic(err)
		}

		configInstance = &cfg
	})

	return  configInstance
}
