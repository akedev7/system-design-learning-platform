package unit

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go-course-service/internal/models"
	"go-course-service/internal/service"
)

func TestDiagramValidationService_ValidateDiagram_ExactMatch(t *testing.T) {
	dvs := service.NewDiagramValidationService()

	expected := models.DiagramConfig{
		NodeTypes: map[string]int{
			"LoadBalancer": 1,
			"AppServer":     2,
			"Database":      1,
		},
		Edges: []models.DiagramEdge{
			{From: "LoadBalancer", To: "AppServer"},
			{From: "AppServer", To: "Database"},
		},
	}

	userDiagram := models.UserDiagram{
		Nodes: []models.DiagramNode{
			{ID: "1", Type: "LoadBalancer"},
			{ID: "2", Type: "AppServer"},
			{ID: "3", Type: "AppServer"},
			{ID: "4", Type: "Database"},
		},
		Edges: []models.DiagramEdge{
			{From: "1", To: "2"},
			{From: "1", To: "3"},
			{From: "2", To: "4"},
			{From: "3", To: "4"},
		},
	}

	result := dvs.ValidateDiagram(expected, userDiagram)

	assert.True(t, result.Valid)
	assert.Equal(t, 100, result.Score)
	assert.Empty(t, result.Errors)
}

func TestDiagramValidationService_ValidateDiagram_WrongNodeCount(t *testing.T) {
	dvs := service.NewDiagramValidationService()

	expected := models.DiagramConfig{
		NodeTypes: map[string]int{
			"LoadBalancer": 1,
			"AppServer":     2,
			"Database":      1,
		},
		Edges: []models.DiagramEdge{
			{From: "LoadBalancer", To: "AppServer"},
		},
	}

	userDiagram := models.UserDiagram{
		Nodes: []models.DiagramNode{
			{ID: "1", Type: "LoadBalancer"},
			{ID: "2", Type: "AppServer"},
			{ID: "3", Type: "AppServer"},
			{ID: "4", Type: "AppServer"}, // 3 instead of 2
		},
		Edges: []models.DiagramEdge{
			{From: "1", To: "2"},
		},
	}

	result := dvs.ValidateDiagram(expected, userDiagram)

	assert.False(t, result.Valid)
	assert.Less(t, result.Score, 100)
	assert.NotEmpty(t, result.Errors)
}

func TestDiagramValidationService_ValidateDiagram_MissingEdge(t *testing.T) {
	dvs := service.NewDiagramValidationService()

	expected := models.DiagramConfig{
		NodeTypes: map[string]int{
			"LoadBalancer": 1,
			"AppServer":     1,
		},
		Edges: []models.DiagramEdge{
			{From: "LoadBalancer", To: "AppServer"},
		},
	}

	userDiagram := models.UserDiagram{
		Nodes: []models.DiagramNode{
			{ID: "1", Type: "LoadBalancer"},
			{ID: "2", Type: "AppServer"},
		},
		Edges: []models.DiagramEdge{
			// Missing connection
		},
	}

	result := dvs.ValidateDiagram(expected, userDiagram)

	assert.False(t, result.Valid)
	assert.NotEmpty(t, result.Errors)
}

func TestDiagramValidationService_ValidateDiagram_MissingNodes(t *testing.T) {
	dvs := service.NewDiagramValidationService()

	expected := models.DiagramConfig{
		NodeTypes: map[string]int{
			"LoadBalancer": 1,
			"AppServer":     1,
			"Database":      1,
		},
		Edges: []models.DiagramEdge{},
	}

	userDiagram := models.UserDiagram{
		Nodes: []models.DiagramNode{
			{ID: "1", Type: "LoadBalancer"},
			{ID: "2", Type: "AppServer"},
			// Missing Database
		},
		Edges: []models.DiagramEdge{},
	}

	result := dvs.ValidateDiagram(expected, userDiagram)

	assert.False(t, result.Valid)
	assert.NotEmpty(t, result.Errors)
}

func TestDiagramValidationService_ValidateDiagram_ExtraNodes(t *testing.T) {
	dvs := service.NewDiagramValidationService()

	expected := models.DiagramConfig{
		NodeTypes: map[string]int{
			"LoadBalancer": 1,
			"AppServer":     1,
		},
		Edges: []models.DiagramEdge{},
	}

	userDiagram := models.UserDiagram{
		Nodes: []models.DiagramNode{
			{ID: "1", Type: "LoadBalancer"},
			{ID: "2", Type: "AppServer"},
			{ID: "3", Type: "Cache"}, // Extra node
		},
		Edges: []models.DiagramEdge{},
	}

	result := dvs.ValidateDiagram(expected, userDiagram)

	assert.False(t, result.Valid)
	assert.NotEmpty(t, result.Errors)
}

func TestDiagramValidationService_ValidateDiagram_ComplexDiagram(t *testing.T) {
	dvs := service.NewDiagramValidationService()

	expected := models.DiagramConfig{
		NodeTypes: map[string]int{
			"Client":        1,
			"APIGateway":    1,
			"LoadBalancer":  1,
			"AppServer":     2,
			"Cache":         1,
			"Database":      1,
			"CDN":           1,
		},
		Edges: []models.DiagramEdge{
			{From: "Client", To: "CDN"},
			{From: "CDN", To: "APIGateway"},
			{From: "APIGateway", To: "LoadBalancer"},
			{From: "LoadBalancer", To: "AppServer"},
			{From: "AppServer", To: "Cache"},
			{From: "AppServer", To: "Database"},
		},
	}

	userDiagram := models.UserDiagram{
		Nodes: []models.DiagramNode{
			{ID: "1", Type: "Client"},
			{ID: "2", Type: "CDN"},
			{ID: "3", Type: "APIGateway"},
			{ID: "4", Type: "LoadBalancer"},
			{ID: "5", Type: "AppServer"},
			{ID: "6", Type: "AppServer"},
			{ID: "7", Type: "Cache"},
			{ID: "8", Type: "Database"},
		},
		Edges: []models.DiagramEdge{
			{From: "1", To: "2"},
			{From: "2", To: "3"},
			{From: "3", To: "4"},
			{From: "4", To: "5"},
			{From: "4", To: "6"},
			{From: "5", To: "7"},
			{From: "6", To: "7"},
			{From: "5", To: "8"},
			{From: "6", To: "8"},
		},
	}

	result := dvs.ValidateDiagram(expected, userDiagram)

	assert.True(t, result.Valid)
	assert.Equal(t, 100, result.Score)
}

func TestDiagramValidationService_ValidateDiagram_EmptyDiagram(t *testing.T) {
	dvs := service.NewDiagramValidationService()

	expected := models.DiagramConfig{
		NodeTypes: map[string]int{
			"AppServer": 1,
		},
		Edges: []models.DiagramEdge{},
	}

	userDiagram := models.UserDiagram{
		Nodes: []models.DiagramNode{},
		Edges: []models.DiagramEdge{},
	}

	result := dvs.ValidateDiagram(expected, userDiagram)

	assert.False(t, result.Valid)
	assert.NotEmpty(t, result.Errors)
}

func TestDiagramValidationService_ValidateDiagram_PartialCorrect(t *testing.T) {
	dvs := service.NewDiagramValidationService()

	expected := models.DiagramConfig{
		NodeTypes: map[string]int{
			"LoadBalancer": 1,
			"AppServer":     1,
			"Database":      1,
		},
		Edges: []models.DiagramEdge{
			{From: "LoadBalancer", To: "AppServer"},
			{From: "AppServer", To: "Database"},
		},
	}

	userDiagram := models.UserDiagram{
		Nodes: []models.DiagramNode{
			{ID: "1", Type: "LoadBalancer"},
			{ID: "2", Type: "AppServer"},
			{ID: "3", Type: "Database"},
		},
		Edges: []models.DiagramEdge{
			{From: "1", To: "2"},
			// Missing AppServer to Database
		},
	}

	result := dvs.ValidateDiagram(expected, userDiagram)

	assert.False(t, result.Valid)
	assert.Less(t, result.Score, 100)
	assert.Greater(t, result.Score, 0)
}