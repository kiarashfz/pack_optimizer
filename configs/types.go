package configs

type Config struct {
	App `mapstructure:",squash"` // Squash allows the fields of App to be directly accessible in Config
	DB  `mapstructure:",squash"`
	Env string `mapstructure:"ENV" validate:"required"`
}

type App struct {
	Port string `mapstructure:"APP_PORT" validate:"required"`
}

type DB struct {
	Host       string `mapstructure:"DB_HOST" validate:"required"`
	Port       string `mapstructure:"DB_PORT" validate:"required"`
	User       string `mapstructure:"DB_USER" validate:"required"`
	Password   string `mapstructure:"DB_PASSWORD" validate:"required"`
	Name       string `mapstructure:"DB_NAME" validate:"required"`
	SSLMode    string `mapstructure:"DB_SSLMODE" validate:"required"`
	GormDSN    string `mapstructure:"-"`
	MigrateDSN string `mapstructure:"-"`
}
