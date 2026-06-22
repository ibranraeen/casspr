package backend

import (
	"github.com/ibranraeen/casspr/internal/proto"
	"github.com/ibranraeen/casspr/internal/question"
)

// SubmitQuestion submits an answer response to a clarifying question.
func (b *Backend) SubmitQuestion(workspaceID string, req proto.QuestionResponse) (bool, error) {
	ws, err := b.GetWorkspace(workspaceID)
	if err != nil {
		return false, err
	}

	resp := question.QuestionResponse{
		QuestionID:      req.QuestionID,
		SelectedOptions: req.SelectedOptions,
		CustomAnswer:    req.CustomAnswer,
	}

	return ws.Questions.Submit(resp), nil
}
