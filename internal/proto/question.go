package proto

// CreateQuestionRequest represents a request to ask a question.
type CreateQuestionRequest struct {
	SessionID     string   `json:"session_id"`
	ToolCallID    string   `json:"tool_call_id"`
	Question      string   `json:"question"`
	Options       []string `json:"options"`
	IsMultiSelect bool     `json:"is_multi_select"`
	AllowCustom   bool     `json:"allow_custom"`
}

// QuestionRequest represents a pending clarifying question.
type QuestionRequest struct {
	ID            string   `json:"id"`
	SessionID     string   `json:"session_id"`
	ToolCallID    string   `json:"tool_call_id"`
	Question      string   `json:"question"`
	Options       []string `json:"options"`
	IsMultiSelect bool     `json:"is_multi_select"`
	AllowCustom   bool     `json:"allow_custom"`
}

// QuestionResponse represents the user's answer to a question.
type QuestionResponse struct {
	QuestionID      string   `json:"question_id"`
	SelectedOptions []string `json:"selected_options"`
	CustomAnswer    string   `json:"custom_answer,omitempty"`
}

// QuestionNotification represents a notification about a question state.
type QuestionNotification struct {
	ToolCallID string            `json:"tool_call_id"`
	Response   *QuestionResponse `json:"response,omitempty"`
}

// QuestionSubmitResponse is the response to a submitted question.
type QuestionSubmitResponse struct {
	Resolved bool `json:"resolved"`
}
