package model

// Person : struct for user
type Person struct {
	ID           string `bson,json:"id"`
	Name         string `bson,json:"name"`
	Works        bool   `bson,json:"works"`
	Age          int32  `bson,json:"age"`
	Password     string `bson,json:"password"`
	RefreshToken string `bson,json:"refreshToken"`
}

// Authentication struct for parse it
/*type Authentication struct {
	Password string `json:"password"`
}

// RefreshTokens struct for parse it
type RefreshTokens struct {
	RefreshToken string `json:"refreshToken"`
}

// Response struct create response
type Response struct {
	Message  string
	FileType string
	FileSize int64
}

// Config struct create config
type Config struct {
	CurrentDB     string `env:"CURRENT_DB" envDefault:"postgres"`
	Password      string `env:"PASSWORD"`
	PostgresDBURL string `env:"POSTGRES_DB_URL"`
	MongoDBURL    string `env:"MONGO_DB_URL"`
	RedisURL      string `env:"REDIS_DB_URL" envDefault:"localhost:6379"`
}
*/
