// Package model models that use in this project
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

// Config struct create config
type Config struct {
	CurrentDB     string `env:"CURRENT_DB" envDefault:"postgres"`
	PostgresDBURL string `env:"POSTGRES_DB_URL"`
	MongoDBURL    string `env:"MONGO_DB_URL"`
	JwtKey        []byte `env:"JWT-KEY" envDefault:"super-key"`
}
