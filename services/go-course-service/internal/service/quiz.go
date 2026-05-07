package service

import "go-course-service/internal/models"

const passingThreshold = 60

type QuizService struct{}

func NewQuizService() *QuizService {
	return &QuizService{}
}

func (s *QuizService) GradeQuiz(quiz models.Quiz, answers map[int]string) models.QuizResult {
	if len(quiz.Questions) == 0 {
		return models.QuizResult{
			Score:          0,
			TotalQuestions: 0,
			CorrectAnswers: 0,
			Passed:         false,
		}
	}

	correctAnswers := 0
	for _, question := range quiz.Questions {
		if answer, ok := answers[question.ID]; ok {
			if answer == question.Correct {
				correctAnswers++
			}
		}
	}

	totalQuestions := len(quiz.Questions)
	score := (correctAnswers * 100) / totalQuestions
	passed := score >= passingThreshold

	return models.QuizResult{
		Score:          score,
		TotalQuestions: totalQuestions,
		CorrectAnswers: correctAnswers,
		Passed:         passed,
	}
}