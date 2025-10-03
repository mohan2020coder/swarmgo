package agent

import (
    "context"
    "encoding/json"
    "fmt"
    "net/http"
    "bytes"
    "io/ioutil"
    "time"

    "github.com/example/document-builder/internal/logging"
)

type styleAgent struct {
    logger *logging.Logger
    serviceURL string
}

func NewStyleAgent(logger *logging.Logger) *styleAgent {
    return &styleAgent{logger: logger, serviceURL: "http://docx-service:5000/generate"}
}

func (s *styleAgent) GenerateDocx(ctx context.Context, req DocxRequest) (string, error) {
    b, _ := json.Marshal(req)
    r, err := http.NewRequestWithContext(ctx, "POST", s.serviceURL, bytes.NewReader(b))
    if err != nil {
        return "", err
    }
    r.Header.Set("Content-Type", "application/json")
    client := &http.Client{Timeout: 90 * time.Second}
    resp, err := client.Do(r)
    if err != nil {
        return "", err
    }
    defer resp.Body.Close()
    data, _ := ioutil.ReadAll(resp.Body)
    if resp.StatusCode >= 400 {
        return "", fmt.Errorf("docx service err: %s", string(data))
    }
    var out struct{ Path string `json:"path"` }
    if err := json.Unmarshal(data, &out); err != nil {
        return "", err
    }
    return out.Path, nil
}
