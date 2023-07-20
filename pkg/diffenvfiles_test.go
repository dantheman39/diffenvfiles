package pkg

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func getEnvFile(t *testing.T, path string) []byte {
	data, err := os.ReadFile(path)
	assert.NoError(t, err)
	return data
}

func getGoodEnv1(t *testing.T) []byte {
	return getEnvFile(t, "./pkg/fixtures/goodEnv1")
}

func getGoodEnv2(t *testing.T) []byte {
	return getEnvFile(t, "./pkg/fixtures/goodEnv2")
}

func TestParseFile(t *testing.T) {
	goodEnv1 := getGoodEnv1(t)
	result, err := parseContents(goodEnv1)
	assert.NoError(t, err)
	envs := result.EnvVars
	assert.Equal(t, 4, len(envs))
	assert.Equal(t, "ABC", envs[0][0])
	assert.Equal(t, "\"hey there\"", envs[0][1])
	assert.Equal(t, "DEF", envs[1][0])
	assert.Equal(t, "\"other=okay yeah\"", envs[1][1])
	assert.Equal(t, "ONLY_IN_ENV_1", envs[2][0])
	assert.Equal(t, "\"only1\"", envs[2][1])
	assert.Equal(t, "SAME_KEY_DIFFERENT_VALUE", envs[3][0])
	assert.Equal(t, "\"env1val\"", envs[3][1])

}
