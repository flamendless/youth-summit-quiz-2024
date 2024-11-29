package handlers

import (
	"fmt"
	"net/http"
	"strings"
	"youth-summit-quiz-2024/client/components"
	"youth-summit-quiz-2024/internal/constants"
	"youth-summit-quiz-2024/internal/models"

	"github.com/a-h/templ"
	"go.uber.org/zap"
)

type HomeService struct {
	QAs []*models.QA
}

func NewHomeService(qas []*models.QA) *HomeService {
	return &HomeService{
		QAs: qas,
	}
}

type HomeHandler struct {
	Logger      *zap.Logger
	HomeService *HomeService
}

func NewHomeHandler(
	logger *zap.Logger,
	homeService *HomeService,
) HomeHandler {
	return HomeHandler{
		Logger:      logger,
		HomeService: homeService,
	}
}

func (h HomeHandler) HomePage(w http.ResponseWriter, r *http.Request) {
	cards := make([]templ.Component, 0, 90-65)
	for i := constants.ASCII_A; i <= constants.ASCII_Z; i++ {
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
	questionToUse := models.GetQuestion(h.HomeService.QAs, letter, difficulty)

	components.Base(
		title,
		components.WideCenterCard(
			components.CardTitle(title, letter),
			components.Question(
				strings.ToUpper(fmt.Sprintf("%s %s", questionToUse.Difficulty, questionToUse.Letter)),
				questionToUse,
			),
		),
	).Render(r.Context(), w)
}

func (h HomeHandler) AnswerPage(w http.ResponseWriter, r *http.Request) {
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
	questionToUse := models.GetQuestion(h.HomeService.QAs, letter, difficulty)

	components.Base(
		title,
		components.WideCenterCard(
			components.CardTitle(title, letter),
			components.Answer(
				strings.ToUpper(fmt.Sprintf("%s %s", questionToUse.Difficulty, questionToUse.Letter)),
				questionToUse,
			),
		),
	).Render(r.Context(), w)
}
