package logger

type Env struct {
	Level     string `env:"LOG_LEVEL" envDefault:""`
	Formatter string `env:"LOG_FORMATTER" envDefault:"json"`
}
