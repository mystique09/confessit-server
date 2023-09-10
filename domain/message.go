package domain

import (
	"time"

	"github.com/google/uuid"
)

type (
	IMessageID = IBaseField

	MessageID struct {
		value uuid.UUID
	}

	IMessage interface {
		ID() IMessageID
		ReceiverID() IUserID
		Content() string
		Seen() bool
		IDateFields
	}

	Message struct {
		id          IMessageID
		receiver_id IUserID
		content     string
		seen        bool
		created_at  time.Time
		updated_at  time.Time
	}

	CreateMessageDTO struct {
		ReceiverID uuid.UUID `json:"receiver_id"`
		Content    string    `json:"content"`
	}
)

func NewMessageID() IMessageID {
	return MessageID{value: uuid.New()}
}

func (m_id MessageID) String() string {
	return m_id.value.String()
}

func (m_id MessageID) ValidateLength(n int) bool {
	return len(m_id.String()) == n
}

func NewMessage(receiver_id IUserID, content string) IMessage {
	return Message{
		id:          NewMessageID(),
		receiver_id: receiver_id,
		content:     content,
		seen:        false,
		created_at:  time.Now(),
		updated_at:  time.Now(),
	}
}

func (m Message) ID() IMessageID {
	return m.id
}

func (m Message) ReceiverID() IUserID {
	return m.receiver_id
}

func (m Message) Content() string {
	return m.content
}

func (m Message) Seen() bool {
	return m.seen
}

func (m Message) CreatedAt() time.Time {
	return m.created_at
}

func (m Message) UpdatedAt() time.Time {
	return m.updated_at
}

type MessageResponse struct {
	ID         IMessageID `json:"id"`
	ReceiverID IUserID    `json:"receiver_id"`
	Content    string     `json:"content"`
	Seen       bool       `json:"seen"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
}

func (m Message) IntoResponse() Response[MessageResponse] {
	return Response[MessageResponse]{
		Message: "",
		Data: MessageResponse{
			ID:         m.ID(),
			ReceiverID: m.ReceiverID(),
			Content:    m.Content(),
			Seen:       m.Seen(),
			CreatedAt:  m.CreatedAt(),
			UpdatedAt:  m.UpdatedAt(),
		},
	}
}

func (payload CreateMessageDTO) IntoMessage() IMessage {
	return NewMessage(UserID.From(UserID{}, payload.ReceiverID), payload.Content)
}
