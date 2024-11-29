package handlers

import (
	"fmt"
	"net/http"
	"strings"
	"youth-summit-quiz-2024/client/components"

	"github.com/a-h/templ"
	"go.uber.org/zap"
)

type HomeService struct {
}

type HomeHandler struct {
	Logger      *zap.Logger
	HomeService HomeService
}

func NewHomeHandler(
	logger *zap.Logger,
	homeService HomeService,
) HomeHandler {
	return HomeHandler{
		Logger:      logger,
		HomeService: homeService,
	}
}

func (h HomeHandler) HomePage(w http.ResponseWriter, r *http.Request) {
	const ASCII_A = 65
	const ASCII_Z = 90

	cards := make([]templ.Component, 0, 90-65)
	for i := ASCII_A; i <= ASCII_Z; i++ {
		char := string(rune(i))
		cards = append(cards, components.LetterCard(char, strings.ToLower(char)))
	}

	components.Base(
		"Youth Summit Quiz 2024",
		components.WideCenterCard(
			cards...,
		),
	).Render(r.Context(), w)
}

func (h HomeHandler) DifficultyPage(w http.ResponseWriter, r *http.Request) {
	letter := r.URL.Query().Get("letter")
	if letter == "" {
		return
	}

	title := "Letter: " + letter
	title = strings.ToUpper(title)

	components.Base(
		title,
		components.WideCenterCard(
			components.CardTitle(title, letter),
			components.DifficultyCardBody(
				components.DifficultyCard(letter, "sprint"),
				components.DifficultyCard(letter, "marathon"),
				components.DifficultyCard(letter, "hurdle"),
			),
		),
	).Render(r.Context(), w)
}

func (h HomeHandler) QuestionPage(w http.ResponseWriter, r *http.Request) {
	letter := r.URL.Query().Get("letter")
	if letter == "" {
		return
	}

	difficulty := r.URL.Query().Get("difficulty")
	if difficulty == "" {
		return
	}
	title := fmt.Sprintf("%s: %s", difficulty, letter)
	title = strings.ToUpper(title)

	components.Base(
		title,
		components.WideCenterCard(
			components.CardTitle(title, letter),
		),
	).Render(r.Context(), w)
}
