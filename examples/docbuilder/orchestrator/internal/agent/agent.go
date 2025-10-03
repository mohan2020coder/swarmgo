package agent

import "context"

// Outline holds list of section titles
type Outline struct {
    Title    string   `json:"title"`
    Sections []string `json:"sections"`
}

// DocxRequest: payload for style/output agent
type DocxRequest struct {
    Title    string            `json:"title"`
    Outline  Outline           `json:"outline"`
    Content  map[string]string `json:"content"`
    Diagrams map[string]string `json:"diagrams"`
    Meta     map[string]string `json:"meta"`
}

type ContentAgent interface {
    GenerateOutline(ctx context.Context, topic, audience, level string) (Outline, error)
    GenerateSection(ctx context.Context, topic, section string) (string, error)
}

type DiagramAgent interface {
    GenerateDiagram(ctx context.Context, section string) (string, error)
}

type StyleAgent interface {
    GenerateDocx(ctx context.Context, req DocxRequest) (string, error)
}
