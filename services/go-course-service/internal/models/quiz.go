package models

import "encoding/json"

type Quiz struct {
	Questions []Question `json:"questions"`
}

type Question struct {
	ID      int      `json:"id"`
	Type    string   `json:"type"`
	Options []string `json:"options,omitempty"`
	Correct string   `json:"correct"`
}

type QuizAnswerRequest struct {
	Answers map[int]string `json:"answers"`
}

type QuizResult struct {
	Score          int  `json:"score"`
	TotalQuestions int  `json:"totalQuestions"`
	CorrectAnswers int  `json:"correctAnswers"`
	Passed         bool `json:"passed"`
}

type ContentBlock struct {
	Type   string          `json:"type"`
	Order  int             `json:"order"`
	Config json.RawMessage `json:"config"`
}

type DiagramConfig struct {
	NodeTypes map[string]int        `json:"nodeTypes"`
	Edges     []DiagramEdge         `json:"edges"`
}

type DiagramEdge struct {
	From string `json:"from"`
	To   string `json:"to"`
}

type UserDiagram struct {
	Nodes []DiagramNode `json:"nodes"`
	Edges []DiagramEdge `json:"edges"`
}

type DiagramNode struct {
	ID   string `json:"id"`
	Type string `json:"type"`
}

type DiagramValidationResult struct {
	Valid  bool     `json:"valid"`
	Score  int      `json:"score"`
	Errors []string `json:"errors"`
}

type DiagramValidationRequest struct {
	Diagram UserDiagram `json:"diagram"`
}