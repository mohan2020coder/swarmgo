package llm

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/sashabaranov/go-openai"
)

// OpenRouterLLM implements the LLM interface for OpenRouter
type OpenRouterLLM struct {
	client *openai.Client
}

// NewOpenRouterLLM creates a new OpenRouter LLM client
func NewOpenRouterLLM(apiKey string) *OpenRouterLLM {
	config := openai.DefaultConfig(apiKey)
	config.BaseURL = "https://openrouter.ai/api/v1"
	config.HTTPClient = &http.Client{} // can customize if needed
	return &OpenRouterLLM{client: openai.NewClientWithConfig(config)}
}

func NewOpenRouterLLMWithHost(apiKey string, host string) *OpenRouterLLM {
	config := openai.DefaultConfig(apiKey)
	if strings.TrimSpace(host) != "" {
		config.BaseURL = host
	} else {
		config.BaseURL = "https://openrouter.ai/api/v1"
	}
	config.HTTPClient = &http.Client{}
	return &OpenRouterLLM{client: openai.NewClientWithConfig(config)}
}

// convertToOpenRouterMessages converts our generic Message type to OpenAI's message type (same format)
func convertToOpenRouterMessages(messages []Message) []openai.ChatCompletionMessage {
	openRouterMessages := []openai.ChatCompletionMessage{}
	for _, msg := range messages {
		if msg.Content == "" {
			continue
		}
		openRouterMessages = append(openRouterMessages, openai.ChatCompletionMessage{
			Role:    string(msg.Role),
			Content: msg.Content,
			Name:    msg.Name,
		})
	}
	return openRouterMessages
}

// convertFromOpenRouterMessage converts OpenRouter's message type to our generic Message type
func convertFromOpenRouterMessage(msg openai.ChatCompletionMessage) Message {
	return Message{
		Role:    Role(msg.Role),
		Content: msg.Content,
		Name:    msg.Name,
	}
}

// convertFromOpenRouterDelta converts OpenRouter's delta message type to our generic Message type
func convertFromOpenRouterDelta(delta openai.ChatCompletionStreamChoiceDelta) Message {
	return Message{
		Role:    Role(delta.Role),
		Content: delta.Content,
	}
}

// convertToOpenRouterTools converts our generic Tool type to OpenRouter's tool type
func convertToOpenRouterTools(tools []Tool) []openai.Tool {
	if len(tools) == 0 {
		return nil
	}

	openRouterTools := make([]openai.Tool, len(tools))
	for i, tool := range tools {
		def := openai.FunctionDefinition{
			Name:        tool.Function.Name,
			Description: tool.Function.Description,
			Parameters:  tool.Function.Parameters,
		}
		openRouterTools[i] = openai.Tool{
			Type:     openai.ToolTypeFunction,
			Function: &def,
		}
	}
	return openRouterTools
}

// convertFromOpenRouterToolCalls converts OpenRouter's tool calls to our generic type
func convertFromOpenRouterToolCalls(toolCalls []openai.ToolCall) []ToolCall {
	if len(toolCalls) == 0 {
		return nil
	}

	calls := make([]ToolCall, len(toolCalls))
	for i, call := range toolCalls {
		calls[i] = ToolCall{
			ID:   call.ID,
			Type: string(call.Type),
		}
		calls[i].Function.Name = call.Function.Name
		calls[i].Function.Arguments = call.Function.Arguments
	}
	return calls
}

// CreateChatCompletion implements the LLM interface for OpenRouter
func (o *OpenRouterLLM) CreateChatCompletion(ctx context.Context, req ChatCompletionRequest) (ChatCompletionResponse, error) {
	openRouterReq := openai.ChatCompletionRequest{
		Model:           req.Model,
		Messages:        convertToOpenRouterMessages(req.Messages),
		Temperature:     float32(req.Temperature),
		TopP:            float32(req.TopP),
		N:               req.N,
		Stop:            req.Stop,
		MaxTokens:       req.MaxTokens,
		PresencePenalty: req.PresencePenalty,
		Tools:           convertToOpenRouterTools(req.Tools),
	}

	resp, err := o.client.CreateChatCompletion(ctx, openRouterReq)
	if err != nil {
		return ChatCompletionResponse{}, err
	}

	choices := make([]Choice, len(resp.Choices))
	for i, c := range resp.Choices {
		msg := convertFromOpenRouterMessage(c.Message)
		msg.ToolCalls = convertFromOpenRouterToolCalls(c.Message.ToolCalls)
		choices[i] = Choice{
			Index:        c.Index,
			Message:      msg,
			FinishReason: string(c.FinishReason),
		}
	}

	return ChatCompletionResponse{
		ID:      resp.ID,
		Choices: choices,
		Usage: Usage{
			PromptTokens:     resp.Usage.PromptTokens,
			CompletionTokens: resp.Usage.CompletionTokens,
			TotalTokens:      resp.Usage.TotalTokens,
		},
	}, nil
}

// openRouterStreamWrapper wraps the OpenRouter stream
type openRouterStreamWrapper struct {
	stream          *openai.ChatCompletionStream
	currentToolCall *ToolCall
	toolCallBuffer  map[string]*ToolCall
}

func newOpenRouterStreamWrapper(stream *openai.ChatCompletionStream) *openRouterStreamWrapper {
	return &openRouterStreamWrapper{
		stream:         stream,
		toolCallBuffer: make(map[string]*ToolCall),
	}
}

func (w *openRouterStreamWrapper) Recv() (ChatCompletionResponse, error) {
	resp, err := w.stream.Recv()
	if err != nil {
		if err == io.EOF {
			return ChatCompletionResponse{}, err
		}
		var apiErr *openai.APIError
		if errors.As(err, &apiErr) {
			return ChatCompletionResponse{}, fmt.Errorf("OpenRouter API error: %s - %s", apiErr.Code, apiErr.Message)
		}
		return ChatCompletionResponse{}, fmt.Errorf("stream receive failed: %w", err)
	}

	choices := make([]Choice, len(resp.Choices))
	for i, c := range resp.Choices {
		message := Message{
			Role:    Role(c.Delta.Role),
			Content: c.Delta.Content,
		}

		// Handle tool calls in delta
		if len(c.Delta.ToolCalls) > 0 {
			message.ToolCalls = make([]ToolCall, 0)
			for _, tc := range c.Delta.ToolCalls {
				toolCall, exists := w.toolCallBuffer[tc.ID]
				if !exists {
					if tc.ID == "" {
						if tc.Function.Arguments != "" && w.currentToolCall != nil {
							w.currentToolCall.Function.Arguments += tc.Function.Arguments
							var args map[string]interface{}
							if err := json.Unmarshal([]byte(w.currentToolCall.Function.Arguments), &args); err == nil {
								message.ToolCalls = append(message.ToolCalls, *w.currentToolCall)
								delete(w.toolCallBuffer, w.currentToolCall.ID)
								w.currentToolCall = nil
							}
						}
						continue
					}
					toolCall = &ToolCall{
						ID:   tc.ID,
						Type: string(tc.Type),
						Function: ToolCallFunction{
							Name:      tc.Function.Name,
							Arguments: "",
						},
					}
					w.toolCallBuffer[tc.ID] = toolCall
					w.currentToolCall = toolCall
				}

				if tc.Function.Name != "" {
					toolCall.Function.Name = tc.Function.Name
				}
				if tc.Function.Arguments != "" {
					toolCall.Function.Arguments += tc.Function.Arguments
					var args map[string]interface{}
					if err := json.Unmarshal([]byte(toolCall.Function.Arguments), &args); err == nil {
						message.ToolCalls = append(message.ToolCalls, *toolCall)
						delete(w.toolCallBuffer, tc.ID)
						w.currentToolCall = nil
					}
				}
			}
		}

		choices[i] = Choice{
			Index:        c.Index,
			Message:      message,
			FinishReason: string(c.FinishReason),
		}
	}

	return ChatCompletionResponse{
		ID:      resp.ID,
		Choices: choices,
	}, nil
}

func (w *openRouterStreamWrapper) Close() error {
	return w.stream.Close()
}

// CreateChatCompletionStream implements the LLM interface for OpenRouter streaming
func (o *OpenRouterLLM) CreateChatCompletionStream(ctx context.Context, req ChatCompletionRequest) (ChatCompletionStream, error) {
	openRouterReq := openai.ChatCompletionRequest{
		Model:           req.Model,
		Messages:        convertToOpenRouterMessages(req.Messages),
		Temperature:     float32(req.Temperature),
		TopP:            float32(req.TopP),
		N:               req.N,
		Stop:            req.Stop,
		MaxTokens:       req.MaxTokens,
		PresencePenalty: float32(req.PresencePenalty),
		Tools:           convertToOpenRouterTools(req.Tools),
		Stream:          true,
	}

	stream, err := o.client.CreateChatCompletionStream(ctx, openRouterReq)
	if err != nil {
		var apiErr *openai.APIError
		if errors.As(err, &apiErr) {
			return nil, fmt.Errorf("OpenRouter API error: %s - %s", apiErr.Code, apiErr.Message)
		}
		return nil, fmt.Errorf("stream creation failed: %w", err)
	}

	return newOpenRouterStreamWrapper(stream), nil
}
