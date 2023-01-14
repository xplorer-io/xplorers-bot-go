package ssm

import (
	"encoding/json"
	"net/http"
	"os"
)

type Ssm struct {
	Parameter struct {
		Name  string `json:"Name"`
		Value string `json:"Value"`
	}
}

func GetSsmParameter(parameterPath string) (string, error) {
	URL := "http://localhost:2773/systemsmanager/parameters/get?withDecryption=true&name=" + parameterPath

	req, _ := http.NewRequest("GET", URL, nil)
	req.Header.Add("X-Aws-Parameters-Secrets-Token", os.Getenv("AWS_SESSION_TOKEN"))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}

	var ssm Ssm
	if err := json.NewDecoder(resp.Body).Decode(&ssm); err != nil {
		return "", err
	}

	return ssm.Parameter.Value, nil
}
