package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"rate_limite/internal/limiter"
	"rate_limite/internal/middleware"
	"rate_limite/internal/storage"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

func main () {
	err := godotenv.Load()
	if err != nil {
		log.Println("Arquivo .env não encontrado, usando variaveis de ambiente")
	}

	redisAddr := os.Getenv("REDIS_ADDR")
	if redisAddr == "" {
		redisAddr = "localhost:6379"
	}

	defaultLimit, _ := strconv.Atoi(os.Getenv("RATE_LIMIT_IP"))
	if defaultLimit == 0 {
		defaultLimit = 10
	} 
	defaultTokenLimit, _ := strconv.Atoi(os.Getenv("RATE_LIMIT_TOKEN"))
	if defaultTokenLimit == 0 {
		defaultTokenLimit = 100
	}
	blockTime, _ := strconv.Atoi(os.Getenv("BLOCK_TIME"))
	if blockTime == 0 {
		blockTime = 300
	}
	window := 30 *time.Second

	redisStorage := storage.NewRedisStorage(redisAddr)

	rl := limiter.NewRateLimiter(redisStorage, defaultLimit, window, time.Duration(blockTime)*time.Second)

	tokenLimits := map[string]int {}
	tokenBlockTimes := map[string]time.Duration{}
	tokenWindows := map[string]time.Duration{}

	config := &middleware.RateLimiterConfig{
		Limiter: rl,
		TokenLimits: tokenLimits,
		TokenBlockTimes: tokenBlockTimes,
		TokenWindows: tokenWindows,
		DefaultLimit: defaultLimit,
		DefaultBlockTime: time.Duration(blockTime) *time.Second,
		DefaultWindow: window,
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Rate Limiter funcionando!")
	})

	handler := middleware.RateLimiterMiddleware(config)(mux)

	port := os.Getenv("SERVER_PORT")
	if port == "" {
		port = "8080"
	}
	fmt.Println("Servidor rodando na porta", port)
	log.Fatal(http.ListenAndServe(":" + port, handler))
}