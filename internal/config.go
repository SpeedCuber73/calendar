package internal

// Config базовый конфиг приложения
type Config struct {
	HTTPListen     string `config:"host"`             // ip и port на котором должен слушать web-сервер
	LogFile        string `config:"logfile,required"` // путь к файлу логов
	LogLevel       string `config:"loglevel"`         // уровень логирования (error / warn / info / debug)
	PgName         string
	PgHost         string
	PgPort         int
	PgUser         string
	PgPassword     string
	RabbitHost     string
	RabbitPort     int
	RabbitUser     string
	RabbitPassword string
}
