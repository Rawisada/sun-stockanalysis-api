package configurations

import (
	"strings"
	"sync"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

type (
	Config struct {
		Server *Server `mapstructure:"server" validate:"required"`
		// OAuth2   	*OAuth2   	`mapstructure:"oauth2" validate:"required"`
		State    *State    `mapstructure:"state" validate:"required"`
		Database *Database `mapstructure:"database" validate:"required"`
		Finnhub  *Finnhub  `mapstructure:"finnhub" validate:"required"`
		Push     *Push     `mapstructure:"push"`
	}

	Server struct {
		Host           string        `mapstructure:"host" validate:"required"`
		Port           int           `mapstructure:"port" validate:"required"`
		ContextPath    string        `mapstructure:"contextPath" validate:"required"`
		AllowedOrigins []string      `mapstructure:"allowOrigins" validate:"required"`
		BodyLimit      string        `mapstructure:"bodyLimit" validate:"required"`
		TimeOut        time.Duration `mapstructure:"timeout" validate:"required"`
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

	State struct {
		Secret     string        `mapstructure:"secret" validate:"required"`
		ExpiredsAt time.Duration `mapstructure:"expiredsAt" validate:"required"`
		Issuer     string        `mapstructure:"issuer" validate:"required"`
	}

	Database struct {
		Host     string `mapstructure:"host" validate:"required"`
		Port     int    `mapstructure:"port" validate:"required"`
		User     string `mapstructure:"user" validate:"required"`
		Password string `mapstructure:"password" validate:"required"`
		DBname   string `mapstructure:"dbname" validate:"required"`
		SSLmode  string `mapstructure:"sslmode" validate:"required"`
		Schema   string `mapstructure:"schema" validate:"required"`
	}

	Finnhub struct {
		Token string `mapstructure:"token" validate:"required"`
	}

	Push struct {
		Subject         string `mapstructure:"subject"`
		TriggerScore    int    `mapstructure:"triggerScore"`
		VAPIDPublicKey  string `mapstructure:"vapidPublicKey"`
		VAPIDPrivateKey string `mapstructure:"vapidPrivateKey"`
	}
)

var (
	once           sync.Once //Singleton
	configInstance *Config
)

func ConfigGetting() *Config {
	once.Do(func() {
		_ = godotenv.Load()

		viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
		viper.AutomaticEnv()
		bindEnvKeys()

		if viper.IsSet("server.allowOrigins") {
			raw := strings.TrimSpace(viper.GetString("server.allowOrigins"))
			if raw != "" {
				parts := strings.Split(raw, ",")
				for i := range parts {
					parts[i] = strings.TrimSpace(parts[i])
				}
				viper.Set("server.allowOrigins", parts)
			}
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

	return configInstance
}

func bindEnvKeys() {
	keys := []string{
		"server.host",
		"server.port",
		"server.contextPath",
		"server.allowOrigins",
		"server.bodyLimit",
		"server.timeout",
		"state.secret",
		"state.expiredsAt",
		"state.issuer",
		"database.host",
		"database.port",
		"database.user",
		"database.password",
		"database.dbname",
		"database.sslmode",
		"database.schema",
		"finnhub.token",
		"push.subject",
		"push.triggerScore",
		"push.vapidPublicKey",
		"push.vapidPrivateKey",
	}

	for _, key := range keys {
		_ = viper.BindEnv(key)
	}
}
