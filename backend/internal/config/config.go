package config

import (
	"time"

	"github.com/caarlos0/env/v9"
	_ "github.com/joho/godotenv/autoload"
)

var Conf = struct {
	Debug    bool   `env:"DEBUG" envDefault:"false"`
	LogLevel string `env:"LOG_LEVEL" envDefault:"info"`

	GrpcPort string `env:"GRPC_PORT" envDefault:"5050"`
	HttpPort string `env:"HTTP_PORT" envDefault:"8090"`
	// healthcheck and swagger
	SystemHttpPort string `env:"SYSTEM_HTTP_PORT" envDefault:"3003"`
	HttpCors       bool   `env:"HTTP_CORS" envDefault:"false"`

	PgDsn string `env:"PG_DSN"`

	// public app URL: the Google redirect_uri is built from it
	PublicUrl string `env:"PUBLIC_URL" envDefault:"http://localhost:8090"`
	// built frontend (web/dist); empty means the SPA is not served
	WebDir    string `env:"WEB_DIR" envDefault:"../web/dist"`
	UploadDir string `env:"UPLOAD_DIR" envDefault:"./uploads"`

	JwtSecret  string        `env:"JWT_SECRET" envDefault:"dev-secret-change-me"`
	AccessTtl  time.Duration `env:"ACCESS_TTL" envDefault:"15m"`
	RefreshTtl time.Duration `env:"REFRESH_TTL" envDefault:"720h"`

	AllowedDomains []string `env:"ALLOWED_DOMAINS" envDefault:"nu.edu.kz" envSeparator:","`
	Campuses       []string `env:"CAMPUSES" envDefault:"academic_block,dormitory,family_housing,off_campus" envSeparator:","`

	AuthDevLogin       bool   `env:"AUTH_DEV_LOGIN" envDefault:"false"`
	GoogleClientId     string `env:"GOOGLE_CLIENT_ID"`
	GoogleClientSecret string `env:"GOOGLE_CLIENT_SECRET"`
}{}

func init() {
	if err := env.Parse(&Conf); err != nil {
		panic(err)
	}
}
