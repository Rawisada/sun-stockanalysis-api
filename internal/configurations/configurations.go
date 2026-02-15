package configurations

import (
	"bytes"
	"os"
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
		loadDotEnv()

		viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
		viper.AutomaticEnv()
		bindEnvKeys()

		if viper.IsSet("server.allowOrigins") {
			raw := strings.TrimSpace(viper.GetString("server.allowOrigins"))
			if raw != "" {
				parts := strings.Split(raw, ",")
				cleaned := make([]string, 0, len(parts))
				for i := range parts {
					parts[i] = strings.TrimSpace(parts[i])
					if parts[i] != "" {
						cleaned = append(cleaned, parts[i])
					}
				}
				if len(cleaned) > 0 {
					viper.Set("server.allowOrigins", cleaned)
				}
			}
		}

		cfg := Config{
			Server: &Server{
				Host:           viper.GetString("server.host"),
				Port:           viper.GetInt("server.port"),
				ContextPath:    viper.GetString("server.contextPath"),
				AllowedOrigins: viper.GetStringSlice("server.allowOrigins"),
				BodyLimit:      viper.GetString("server.bodyLimit"),
				TimeOut:        viper.GetDuration("server.timeout"),
			},
			State: &State{
				Secret:     viper.GetString("state.secret"),
				ExpiredsAt: viper.GetDuration("state.expiredsAt"),
				Issuer:     viper.GetString("state.issuer"),
			},
			Database: &Database{
				Host:     viper.GetString("database.host"),
				Port:     viper.GetInt("database.port"),
				User:     viper.GetString("database.user"),
				Password: viper.GetString("database.password"),
				DBname:   viper.GetString("database.dbname"),
				SSLmode:  viper.GetString("database.sslmode"),
				Schema:   viper.GetString("database.schema"),
			},
			Finnhub: &Finnhub{
				Token: viper.GetString("finnhub.token"),
			},
			Push: &Push{
				Subject:         viper.GetString("push.subject"),
				TriggerScore:    viper.GetInt("push.triggerScore"),
				VAPIDPublicKey:  viper.GetString("push.vapidPublicKey"),
				VAPIDPrivateKey: viper.GetString("push.vapidPrivateKey"),
			},
		}

		if err := validator.New().Struct(&cfg); err != nil {
			panic(err)
		}

		configInstance = &cfg
	})

	return configInstance
}

func loadDotEnv() {
	paths := []string{
		".env",
		"../.env",
		"../../.env",
		".env.prod",
		"../.env.prod",
		"../../.env.prod",
	}

	for _, path := range paths {
		data, err := os.ReadFile(path)
		if err != nil {
			continue
		}

		// Some Windows editors save UTF-8 files with BOM; trim it for dotenv parsing.
		data = bytes.TrimPrefix(data, []byte{0xEF, 0xBB, 0xBF})

		envMap, err := godotenv.Unmarshal(string(data))
		if err != nil {
			continue
		}

		for key, value := range envMap {
			if _, exists := os.LookupEnv(key); !exists {
				_ = os.Setenv(key, value)
			}
		}
		return
	}
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
