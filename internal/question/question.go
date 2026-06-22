package question

import (
	"context"
	"sync"

	"github.com/google/uuid"
	"github.com/ibranraeen/casspr/internal/csync"
	"github.com/ibranraeen/casspr/internal/pubsub"
)

// CreateQuestionRequest represents options to create a question.
type CreateQuestionRequest struct {
	SessionID     string   `json:"session_id"`
	ToolCallID    string   `json:"tool_call_id"`
	Question      string   `json:"question"`
	Options       []string `json:"options"`
	IsMultiSelect bool     `json:"is_multi_select"`
	AllowCustom   bool     `json:"allow_custom"`
}

// QuestionRequest represents a pending clarifying question request.
type QuestionRequest struct {
	ID            string   `json:"id"`
	SessionID     string   `json:"session_id"`
	ToolCallID    string   `json:"tool_call_id"`
	Question      string   `json:"question"`
	Options       []string `json:"options"`
	IsMultiSelect bool     `json:"is_multi_select"`
	AllowCustom   bool     `json:"allow_custom"`
}

// QuestionResponse represents the user's answer to a clarifying question.
type QuestionResponse struct {
	QuestionID      string   `json:"question_id"`
	SelectedOptions []string `json:"selected_options"`
	CustomAnswer    string   `json:"custom_answer,omitempty"`
}

// QuestionNotification represents a status change notification for a question.
type QuestionNotification struct {
	ToolCallID string            `json:"tool_call_id"`
	Response   *QuestionResponse `json:"response,omitempty"`
}

// Service provides a pub/sub request-response broker for agent questions.
type Service interface {
	pubsub.Subscriber[QuestionRequest]
	Request(ctx context.Context, opts CreateQuestionRequest) (QuestionResponse, error)
	Submit(resp QuestionResponse) bool
	SubscribeNotifications(ctx context.Context) <-chan pubsub.Event[QuestionNotification]
}

type questionService struct {
	*pubsub.Broker[QuestionRequest]
	notificationBroker *pubsub.Broker[QuestionNotification]
	pendingRequests    *csync.Map[string, chan QuestionResponse]

	activeRequest   *QuestionRequest
	activeRequestMu sync.Mutex
}

// Request registers a question request, broadcasts it, and blocks until a
// response is received or the context is cancelled.
func (s *questionService) Request(ctx context.Context, opts CreateQuestionRequest) (QuestionResponse, error) {
	req := QuestionRequest{
		ID:            uuid.New().String(),
		SessionID:     opts.SessionID,
		ToolCallID:    opts.ToolCallID,
		Question:      opts.Question,
		Options:       opts.Options,
		IsMultiSelect: opts.IsMultiSelect,
		AllowCustom:   opts.AllowCustom,
	}

	s.activeRequestMu.Lock()
	s.activeRequest = &req
	s.activeRequestMu.Unlock()

	s.notificationBroker.Publish(pubsub.CreatedEvent, QuestionNotification{
		ToolCallID: req.ToolCallID,
	})

	respCh := make(chan QuestionResponse, 1)
	s.pendingRequests.Set(req.ID, respCh)
	defer s.pendingRequests.Del(req.ID)

	s.Publish(pubsub.CreatedEvent, req)

	select {
	case <-ctx.Done():
		return QuestionResponse{}, ctx.Err()
	case resp := <-respCh:
		return resp, nil
	}
}

// Submit resolves a pending question request with the user's answer.
func (s *questionService) Submit(resp QuestionResponse) bool {
	respCh, ok := s.pendingRequests.Take(resp.QuestionID)
	if !ok {
		return false
	}

	s.activeRequestMu.Lock()
	if s.activeRequest != nil && s.activeRequest.ID == resp.QuestionID {
		s.notificationBroker.Publish(pubsub.CreatedEvent, QuestionNotification{
			ToolCallID: s.activeRequest.ToolCallID,
			Response:   &resp,
		})
		s.activeRequest = nil
	}
	s.activeRequestMu.Unlock()

	respCh <- resp
	return true
}

// SubscribeNotifications registers a listener for question state updates.
func (s *questionService) SubscribeNotifications(ctx context.Context) <-chan pubsub.Event[QuestionNotification] {
	return s.notificationBroker.Subscribe(ctx)
}

// NewQuestionService returns a new clarifying question service.
func NewQuestionService() Service {
	return &questionService{
		Broker:             pubsub.NewBroker[QuestionRequest](),
		notificationBroker: pubsub.NewBroker[QuestionNotification](),
		pendingRequests:    csync.NewMap[string, chan QuestionResponse](),
	}
}
