package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
)

type Server struct {
	redis *redis.Client
}

func init() {
	viper.SetDefault("redis.host", "localhost")
	viper.SetDefault("redis.port", "6773")
	viper.SetDefault("redis.password", "")
	viper.SetDefault("redis.db", "0")

	viper.SetConfigName("config")
	viper.SetConfigType("toml")
	viper.AddConfigPath("/app-dev")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("fatal error config file: %v", err))
	}
}

func (s *Server) Run() {
	fmt.Println("Running...")

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			name := r.URL.Query().Get("name")
			err := s.redis.Set(context.Background(), "name", name, time.Hour.Abs())
			if err != nil {
				log.Println("ERROR: ", err)
			}

		case http.MethodGet:
			name, err := s.redis.Get(context.Background(), "name").Result()
			if err != nil {
				log.Println("ERROR: ", err)
			}
			w.Write([]byte(name))

		}
	})

	serverAdress := fmt.Sprintf(":%s", viper.GetString("server.port"))
	log.Fatal(http.ListenAndServe(serverAdress, nil))
}

func main() {

	redisClient := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", viper.GetString("redis.host"), viper.GetString("redis.port")),
		Password: viper.GetString("redis.password"),
		DB:       viper.GetInt("redis.db"),
	})
	server := Server{redis: redisClient}
	server.Run()
}
