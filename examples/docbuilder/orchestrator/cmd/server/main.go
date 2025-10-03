package main

import (
    "context"
    "encoding/json"
    "log"
    "net/http"
    "time"

    "github.com/example/document-builder/internal/agent"
    "github.com/example/document-builder/internal/cache"
    "github.com/example/document-builder/internal/logging"
)

func main() {
    logger := logging.NewLogger("orchestrator.log")
    c, err := cache.NewSQLiteCache("cache/cache.db")
    if err != nil {
        logger.Fatalf("cache init: %v", err)
    }
    contentAgent := agent.NewContentAgent(c, logger)
    diagramAgent := agent.NewDiagramAgent(c, logger)
    styleAgent := agent.NewStyleAgent(logger)

    mux := http.NewServeMux()
    mux.HandleFunc("/generate", func(w http.ResponseWriter, r *http.Request) {
        var req struct {
            Topic           string `json:"topic"`
            Format          string `json:"format"`
            IncludeDiagrams bool   `json:"include_diagrams"`
            Audience        string `json:"audience"`
            DetailLevel     string `json:"detail_level"`
        }
        if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
            http.Error(w, "bad request", http.StatusBadRequest)
            return
        }

        ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
        defer cancel()

        outline, err := contentAgent.GenerateOutline(ctx, req.Topic, req.Audience, req.DetailLevel)
        if err != nil {
            logger.Errorf("outline: %v", err)
            http.Error(w, "outline error", http.StatusInternalServerError)
            return
        }

        expanded := map[string]string{}
        for _, section := range outline.Sections {
            secText, err := contentAgent.GenerateSection(ctx, req.Topic, section)
            if err != nil {
                logger.Errorf("section %v: %v", section, err)
                continue
            }
            expanded[section] = secText
        }

        diagrams := map[string]string{}
        if req.IncludeDiagrams {
            for _, sec := range outline.Sections {
                d, err := diagramAgent.GenerateDiagram(ctx, sec)
                if err == nil && d != "" {
                    diagrams[sec] = d
                }
            }
        }

        docReq := agent.DocxRequest{
            Title:    "Document Builder Project Plan - " + req.Topic,
            Outline:  outline,
            Content:  expanded,
            Diagrams: diagrams,
            Meta:     map[string]string{"audience": req.Audience, "format": req.Format},
        }

        docPath, err := styleAgent.GenerateDocx(ctx, docReq)
        if err != nil {
            logger.Errorf("docx generation: %v", err)
            http.Error(w, "docx error", http.StatusInternalServerError)
            return
        }

        res := map[string]string{"docx_path": docPath}
        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(res)
    })

    srv := &http.Server{
        Addr:         ":8080",
        Handler:      mux,
        ReadTimeout:  15 * time.Second,
        WriteTimeout: 120 * time.Second,
    }

    log.Printf("server listening on %s", srv.Addr)
    if err := srv.ListenAndServe(); err != nil {
        logger.Fatalf("server: %v", err)
    }
}
