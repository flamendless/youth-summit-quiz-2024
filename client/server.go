package client

import (
	"embed"
	"fmt"
	"net/http"
	"youth-summit-quiz-2024/client/handlers"
	"youth-summit-quiz-2024/client/middlewares"
	"youth-summit-quiz-2024/internal/ctx"
	"youth-summit-quiz-2024/internal/logs"

	"go.uber.org/zap"
)

//go:embed static/*
var static embed.FS

func Serve(ctxClient *ctx.ClientFlags) {
	addr := fmt.Sprintf("%s%s", ctxClient.Address, ctxClient.Port)
	logs.Log().Info("Starting site server", zap.String("address", addr))

	logger := logs.Log()

	homeService := handlers.HomeService{}

	homeHandler := handlers.NewHomeHandler(logger, homeService)

	mux := http.NewServeMux()
	mux.Handle("GET /static/", http.FileServer(http.FS(static)))

	//SHOP
	mux.HandleFunc("GET /", homeHandler.HomePage)
	mux.HandleFunc("GET /home", homeHandler.HomePage)
	mux.HandleFunc("GET /difficulty", homeHandler.DifficultyPage)
	mux.HandleFunc("GET /question", homeHandler.QuestionPage)

	mw := middlewares.NewMiddleware(
		mux,
		middlewares.WithSecure(ctxClient.Secure),
		middlewares.WithHTTPOnly(false),
		middlewares.WithRequestDurMetrics(true),
	)

	err := http.ListenAndServe(addr, mw)
	if err != nil {
		panic(err)
	}
}
