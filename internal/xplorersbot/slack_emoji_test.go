package xplorersbot

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetKeyWordsMap(t *testing.T) {
	assert := require.New(t)
	cmd := GetKeyWordsMap()

	assert.NotEmpty(cmd)
	keywords, found := cmd["kubernetes"]
	assert.Equal(keywords, []string{"kube", "kubernetes", "k8s"})
	assert.True(found)
}

func TestGetEmojis(t *testing.T) {
	assert := require.New(t)
	tests := []struct {
		name     string
		text     string
		expected []string
		setup    func()
	}{
		{
			name:     "text has a keyword match for an emoji",
			text:     "What a wondeful app deployed in a Kubernetes environment.",
			expected: []string{"kubernetes"},
			setup: func() {
				GetKeyWordsMap = func() (keywordsMap map[string][]string) {
					return map[string][]string{
						"kubernetes": {"kube", "kubernetes", "k8s"},
						"docker":     {"docker", "containers", "container"},
					}
				}

				GetEmojiMatches = func(emoji string, keywords []string, text string, emojis []string) []string {
					return []string{"kubernetes"}
				}
			},
		},
		{
			name:     "no keyword match found for an emoji",
			text:     "What a wondeful app deployed in a cloud environment.",
			expected: []string(nil),
			setup: func() {
				GetKeyWordsMap = func() (keywordsMap map[string][]string) {
					return map[string][]string{
						"kubernetes": {"kube", "kubernetes", "k8s"},
						"docker":     {"docker", "containers", "container"},
					}
				}

				GetEmojiMatches = func(emoji string, keywords []string, text string, emojis []string) []string {
					return []string(nil)
				}
			},
		},
	}
	oriGetKeyWordsMap := GetKeyWordsMap
	oriGetEmojiMatches := GetEmojiMatches
	for _, test_case := range tests {
		t.Run(test_case.name, func(t *testing.T) {
			result := GetEmojis(test_case.text)
			assert.Equal(test_case.expected, result)
		})
	}
	GetKeyWordsMap = oriGetKeyWordsMap
	GetEmojiMatches = oriGetEmojiMatches
}

func TestGetEmojiMatches(t *testing.T) {
	assert := require.New(t)
	tests := []struct {
		name     string
		emoji    string
		keywords []string
		text     string
		emojis   []string
		expected []string
	}{
		{
			name:     "text has a keyword match for an emoji",
			emoji:    "docker",
			keywords: []string{"docker", "containers", "container"},
			text:     "A wondeful app deployed in a Docker environment.",
			emojis:   []string{},
			expected: []string{"docker"},
		},
		{
			name:     "text has 2 keyword matches for an emoji - should return unique item",
			emoji:    "docker",
			keywords: []string{"docker", "containers", "container"},
			text:     "A wondeful app containerised in a docker environment.",
			emojis:   []string{},
			expected: []string{"docker"},
		},
		{
			name:     "no keyword match found - empty slice",
			emoji:    "docker",
			keywords: []string{"docker", "containers", "container"},
			text:     "Not a matching text is it.",
			emojis:   []string{},
			expected: []string{},
		},
	}

	for _, test_case := range tests {
		t.Run(test_case.name, func(t *testing.T) {
			result := GetEmojiMatches(test_case.emoji, test_case.keywords, test_case.text, test_case.emojis)
			assert.Equal(test_case.expected, result)
		})
	}
}
