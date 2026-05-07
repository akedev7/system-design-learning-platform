package unit

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go-course-service/internal/models"
	"go-course-service/internal/service"
)

func TestQuizService_GradeQuiz_AllCorrect(t *testing.T) {
	qs := service.NewQuizService()

	quiz := models.Quiz{
		Questions: []models.Question{
			{
				ID:       1,
				Type:     "multiple_choice",
				Options:  []string{"A", "B", "C", "D"},
				Correct:  "B",
			},
			{
				ID:       2,
				Type:     "true_false",
				Correct:  "true",
			},
		},
	}

	answers := map[int]string{
		1: "B",
		2: "true",
	}

	result := qs.GradeQuiz(quiz, answers)

	assert.Equal(t, 100, result.Score)
	assert.Equal(t, 2, result.TotalQuestions)
	assert.Equal(t, 2, result.CorrectAnswers)
	assert.True(t, result.Passed)
}

func TestQuizService_GradeQuiz_SomeCorrect(t *testing.T) {
	qs := service.NewQuizService()

	quiz := models.Quiz{
		Questions: []models.Question{
			{
				ID:       1,
				Type:     "multiple_choice",
				Options:  []string{"A", "B", "C", "D"},
				Correct:  "B",
			},
			{
				ID:       2,
				Type:     "true_false",
				Correct:  "true",
			},
		},
	}

	answers := map[int]string{
		1: "B",
		2: "false",
	}

	result := qs.GradeQuiz(quiz, answers)

	assert.Equal(t, 50, result.Score)
	assert.Equal(t, 2, result.TotalQuestions)
	assert.Equal(t, 1, result.CorrectAnswers)
	assert.False(t, result.Passed)
}

func TestQuizService_GradeQuiz_AllWrong(t *testing.T) {
	qs := service.NewQuizService()

	quiz := models.Quiz{
		Questions: []models.Question{
			{
				ID:       1,
				Type:     "multiple_choice",
				Options:  []string{"A", "B", "C", "D"},
				Correct:  "B",
			},
			{
				ID:       2,
				Type:     "true_false",
				Correct:  "true",
			},
		},
	}

	answers := map[int]string{
		1: "C",
		2: "false",
	}

	result := qs.GradeQuiz(quiz, answers)

	assert.Equal(t, 0, result.Score)
	assert.Equal(t, 2, result.TotalQuestions)
	assert.Equal(t, 0, result.CorrectAnswers)
	assert.False(t, result.Passed)
}

func TestQuizService_GradeQuiz_MissingAnswers(t *testing.T) {
	qs := service.NewQuizService()

	quiz := models.Quiz{
		Questions: []models.Question{
			{
				ID:       1,
				Type:     "multiple_choice",
				Options:  []string{"A", "B", "C", "D"},
				Correct:  "B",
			},
			{
				ID:       2,
				Type:     "true_false",
				Correct:  "true",
			},
		},
	}

	answers := map[int]string{
		1: "B",
	}

	result := qs.GradeQuiz(quiz, answers)

	assert.Equal(t, 50, result.Score)
	assert.Equal(t, 2, result.TotalQuestions)
	assert.Equal(t, 1, result.CorrectAnswers)
}

func TestQuizService_GradeQuiz_PassingThreshold(t *testing.T) {
	qs := service.NewQuizService()

	quiz := models.Quiz{
		Questions: []models.Question{
			{
				ID:       1,
				Type:     "multiple_choice",
				Correct:  "A",
			},
			{
				ID:       2,
				Type:     "multiple_choice",
				Correct:  "B",
			},
			{
				ID:       3,
				Type:     "multiple_choice",
				Correct:  "C",
			},
		},
	}

	answers := map[int]string{
		1: "A",
		2: "B",
		3: "D",
	}

	result := qs.GradeQuiz(quiz, answers)

	assert.Equal(t, 66, result.Score)
	assert.True(t, result.Passed)
}