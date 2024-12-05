package client

import (
	"embed"
	"fmt"
	"net/http"
	"youth-summit-quiz-2024/client/handlers"
	"youth-summit-quiz-2024/client/middlewares"
	"youth-summit-quiz-2024/internal/ctx"
	"youth-summit-quiz-2024/internal/logs"
	"youth-summit-quiz-2024/internal/models"

	"go.uber.org/zap"
)

//go:embed static/*
var static embed.FS

func cacheMode(ctxClient *ctx.ClientFlags, next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		if !ctxClient.DevMode {
			return next
		}

		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Cache-Control", "no-store")
			next.ServeHTTP(w, r)
		})
	}(next)
}

func Serve(ctxClient *ctx.ClientFlags) {
	addr := fmt.Sprintf("%s%s", ctxClient.Address, ctxClient.Port)
	logs.Log().Info("Starting site server", zap.String("address", addr))

	logger := logs.Log()

	qas := models.QAsFromMarkdown("./data/questions.md")
	homeService := handlers.NewHomeService(qas)

	homeHandler := handlers.NewHomeHandler(logger, homeService)

	mux := http.NewServeMux()
	mux.Handle(
		"GET /youth-summit-2024-quiz/static/",
		http.StripPrefix(
			"/youth-summit-2024-quiz/",
			cacheMode(
				ctxClient,
				http.FileServer(http.FS(static)),
			),
		),
	)

	//SHOP
	mux.HandleFunc("GET /youth-summit-2024-quiz/", homeHandler.HomePage)
	mux.HandleFunc("GET /youth-summit-2024-quiz/difficulty", homeHandler.DifficultyPage)
	mux.HandleFunc("GET /youth-summit-2024-quiz/question", homeHandler.QuestionPage)
	mux.HandleFunc("GET /youth-summit-2024-quiz/answer", homeHandler.AnswerPage)

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
