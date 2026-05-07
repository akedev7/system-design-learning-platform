package response

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSuccess(t *testing.T) {
	data := map[string]string{"id": "1", "title": "Test Course"}
	env := Success(data)

	assert.Equal(t, StatusSuccess, env.Status)
	assert.Equal(t, data, env.Data)
	assert.Empty(t, env.Message)
}

func TestError(t *testing.T) {
	env := Error("not found")

	assert.Equal(t, StatusError, env.Status)
	assert.Equal(t, "not found", env.Message)
	assert.Nil(t, env.Data)
}

func TestSuccessWithNilData(t *testing.T) {
	env := Success(nil)

	assert.Equal(t, StatusSuccess, env.Status)
	assert.Nil(t, env.Data)
	assert.Empty(t, env.Message)
}