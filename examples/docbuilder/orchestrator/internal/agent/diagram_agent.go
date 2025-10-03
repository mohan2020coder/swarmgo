package agent

import (
    "context"
    "fmt"
    "github.com/example/document-builder/internal/cache"
    "github.com/example/document-builder/internal/logging"
)

type diagramAgent struct {
    cache  *cache.SQLiteCache
    logger *logging.Logger
}

func NewDiagramAgent(c *cache.SQLiteCache, logger *logging.Logger) *diagramAgent {
    return &diagramAgent{cache: c, logger: logger}
}

func (d *diagramAgent) GenerateDiagram(ctx context.Context, section string) (string, error) {
    // For demo: return a simple mermaid flowchart as string
    key := "diagram:" + section
    if v, ok := d.cache.Get(key); ok {
        return v, nil
    }
    mermaid := fmt.Sprintf("flowchart TD\n    A[%s] --> B[Details]", section)
    d.cache.Set(key, mermaid)
    return mermaid, nil
}
