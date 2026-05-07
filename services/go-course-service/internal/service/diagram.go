package service

import (
	"sort"

	"go-course-service/internal/models"
)

type DiagramValidationService struct{}

func NewDiagramValidationService() *DiagramValidationService {
	return &DiagramValidationService{}
}

func (s *DiagramValidationService) ValidateDiagram(expected models.DiagramConfig, userDiagram models.UserDiagram) models.DiagramValidationResult {
	errors := []string{}

	nodeCount := s.countNodeTypes(userDiagram.Nodes)
	errors = append(errors, s.validateNodeCounts(expected.NodeTypes, nodeCount)...)

	edgeErrors := s.validateEdges(expected.Edges, userDiagram.Edges, userDiagram.Nodes)
	errors = append(errors, edgeErrors...)

	score := s.calculateScore(expected, nodeCount, len(errors) > 0)
	valid := len(errors) == 0

	return models.DiagramValidationResult{
		Valid:  valid,
		Score:  score,
		Errors: errors,
	}
}

func (s *DiagramValidationService) countNodeTypes(nodes []models.DiagramNode) map[string]int {
	counts := make(map[string]int)
	for _, node := range nodes {
		counts[node.Type]++
	}
	return counts
}

func (s *DiagramValidationService) validateNodeCounts(expected map[string]int, actual map[string]int) []string {
	errors := []string{}

	for nodeType, expectedCount := range expected {
		actualCount := actual[nodeType]
		if actualCount != expectedCount {
			if actualCount == 0 {
				errors = append(errors, "missing node type: "+nodeType)
			} else if actualCount < expectedCount {
				errors = append(errors, "not enough "+nodeType+" nodes (expected "+string(rune(expectedCount)+'0')+", got "+string(rune(actualCount)+'0')+")")
			} else {
				errors = append(errors, "too many "+nodeType+" nodes (expected "+string(rune(expectedCount)+'0')+", got "+string(rune(actualCount)+'0')+")")
			}
		}
	}

	for nodeType, actualCount := range actual {
		if _, exists := expected[nodeType]; !exists && actualCount > 0 {
			errors = append(errors, "unexpected node type: "+nodeType)
		}
	}

	return errors
}

func (s *DiagramValidationService) validateEdges(expected []models.DiagramEdge, userEdges []models.DiagramEdge, nodes []models.DiagramNode) []string {
	errors := []string{}

	nodeTypeMap := make(map[string]string)
	for _, node := range nodes {
		nodeTypeMap[node.ID] = node.Type
	}

	expectedEdgeSet := make(map[string]bool)
	for _, edge := range expected {
		key := edge.From + "->" + edge.To
		expectedEdgeSet[key] = true
	}

	userEdgeTypes := make(map[string]int)
	for _, edge := range userEdges {
		fromType := nodeTypeMap[edge.From]
		toType := nodeTypeMap[edge.To]
		if fromType == "" || toType == "" {
			continue
		}
		key := fromType + "->" + toType
		userEdgeTypes[key]++
	}

	for edgeKey := range expectedEdgeSet {
		if userEdgeTypes[edgeKey] == 0 {
			parts := splitEdgeKey(edgeKey)
			errors = append(errors, "missing connection: "+parts[0]+" -> "+parts[1])
		}
	}

	return errors
}

func splitEdgeKey(key string) []string {
	for i := 0; i < len(key); i++ {
		if key[i] == '-' && i+1 < len(key) && key[i+1] == '>' {
			return []string{key[:i], key[i+2:]}
		}
	}
	return []string{key, ""}
}

func (s *DiagramValidationService) calculateScore(expected models.DiagramConfig, nodeCounts map[string]int, hasErrors bool) int {
	if len(expected.NodeTypes) == 0 {
		return 0
	}

	totalExpectedNodes := 0
	for _, count := range expected.NodeTypes {
		totalExpectedNodes += count
	}

	if totalExpectedNodes == 0 {
		return 0
	}

	matchedNodes := 0
	for nodeType, expectedCount := range expected.NodeTypes {
		actualCount := nodeCounts[nodeType]
		if actualCount == expectedCount {
			matchedNodes += expectedCount
		} else if actualCount > 0 {
			matchedNodes += min(actualCount, expectedCount)
		}
	}

	nodeScore := (matchedNodes * 100) / totalExpectedNodes

	if len(expected.Edges) > 0 {
		expectedEdgeSet := make(map[string]bool)
		for _, edge := range expected.Edges {
			expectedEdgeSet[edge.From+"->"+edge.To] = true
		}

		matchedEdges := 0
		for range expectedEdgeSet {
			matchedEdges++
		}
		edgeScore := (matchedEdges * 100) / len(expected.Edges)
		nodeScore = (nodeScore + edgeScore) / 2
	}

	if hasErrors {
		nodeScore = nodeScore / 2
	}

	return nodeScore
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func (s *DiagramValidationService) getNodeTypeIDs(nodes []models.DiagramNode, nodeType string) []string {
	var ids []string
	for _, node := range nodes {
		if node.Type == nodeType {
			ids = append(ids, node.ID)
		}
	}
	sort.Strings(ids)
	return ids
}