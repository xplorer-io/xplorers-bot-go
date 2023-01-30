package slack

import (
	"strings"

	"github.com/xplorer-io/xplorers-bot-go/internal/common"
)

var GetKeyWordsMap = func() (keywordsMap map[string][]string) {
	return map[string][]string{
		"kubernetes": {"kube", "kubernetes", "k8s"},
		"docker":     {"docker", "containers", "container"},
		"aws":        {"aws", "lambda", "ec2", "cloudwatch", "codebuild", "kinesis streams", "redshift", "appsync", "ebs", "elasticsarch", "amazon", "s3"},
		"aww-yeah":   {"congratulations", "daami", "way to go", "hurray", "success", "successfully", "kadaa", "sure, will do"},
		"arab":       {"congratulations", "daami", "way to go", "hurray", "success", "successfully", "kadaa", "sure, will do"},
		"celebrate":  {"congratulations", "daami", "way to go", "hurray", "success", "successfully", "kadaa", "sure, will do"},
		"graph_ql":   {"graphql", "graphene", "api query"},
		"mongo_db":   {"mongo", "mongodb", "db", "database"},
		"python":     {"python3", "python", "programming"},
		"github":     {"github", "git", "version control", "versioning", "source control"},
		"git":        {"github", "git", "version control", "versioning", "source control"},
		"react1":     {"reactjs", "react"},
		"nodejs":     {"nodejs", "node"},
		"javascript": {"javascript", "programming"},
		"security":   {"secure", "security", "threat", "hacking", "hacker", "hackers", "enterprise security", "threat vector"},
		"youtube":    {"youtube"},
		"linkedin":   {"linkedin"},
		"typescript": {"typescript"},
		"google":     {"google"},
	}
}

func GetEmojisToReactWith(text string) (emojis []string) {
	for emoji, keywords := range GetKeyWordsMap() {
		emojis = GetEmojiMatches(emoji, keywords, text, emojis)
	}
	return emojis
}

var GetEmojiMatches = func(emoji string, keywords []string, text string, emojis []string) []string {
	for _, keyword := range keywords {
		if strings.Contains(strings.ToLower(text), strings.ToLower(keyword)) && !common.ArrayContainsItem(emojis, emoji) {
			emojis = append(emojis, emoji)
		}
	}
	return emojis
}
