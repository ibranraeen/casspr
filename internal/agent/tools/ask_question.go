package tools

import (
	"context"
	_ "embed"
	"encoding/json"
	"errors"
	"fmt"
	"html/template"

	"charm.land/fantasy"
	"github.com/ibranraeen/casspr/internal/question"
)

const AskQuestionToolName = "ask_question"

//go:embed ask_question.md.tpl
var askQuestionDescriptionTmpl []byte

var askQuestionDescriptionTpl = template.Must(
	template.New("askQuestionDescription").
		Parse(string(askQuestionDescriptionTmpl)),
)

func askQuestionDescription() string {
	return renderTemplate(askQuestionDescriptionTpl, nil)
}

// AskQuestionParams represents the options for asking a clarifying question.
type AskQuestionParams struct {
	Question      string   `json:"question" description:"The question to ask the user to clarify ambiguity."`
	Options       []string `json:"options" description:"The list of multiple choice options. Leave empty if only free-form text input is needed."`
	IsMultiSelect bool     `json:"is_multi_select,omitempty" description:"Set to true if the user can select multiple options. Default is false."`
	AllowCustom   bool     `json:"allow_custom,omitempty" description:"Set to true if the user can type in a custom text response in addition to or instead of options. Default is false."`
}

// NewAskQuestionTool creates a new clarifying question agent tool.
func NewAskQuestionTool(questions question.Service) fantasy.AgentTool {
	return fantasy.NewAgentTool(
		AskQuestionToolName,
		askQuestionDescription(),
		func(ctx context.Context, params AskQuestionParams, call fantasy.ToolCall) (fantasy.ToolResponse, error) {
			if params.Question == "" {
				return fantasy.NewTextErrorResponse("question parameter is required"), nil
			}

			sessionID := GetSessionFromContext(ctx)
			if sessionID == "" {
				return fantasy.ToolResponse{}, errors.New("session ID is required to ask a question")
			}

			resp, err := questions.Request(ctx, question.CreateQuestionRequest{
				SessionID:     sessionID,
				ToolCallID:    call.ID,
				Question:      params.Question,
				Options:       params.Options,
				IsMultiSelect: params.IsMultiSelect,
				AllowCustom:   params.AllowCustom,
			})
			if err != nil {
				return fantasy.ToolResponse{}, fmt.Errorf("error asking question: %w", err)
			}

			outputBytes, err := json.Marshal(resp)
			if err != nil {
				return fantasy.ToolResponse{}, fmt.Errorf("failed to marshal question response: %w", err)
			}

			return fantasy.NewTextResponse(string(outputBytes)), nil
		},
	)
}
