package queue

import "github.com/stretchr/testify/mock"

// EventPublisherMock adalah struct mock untuk EventPublisher
type EventPublisherMock struct {
	mock.Mock
}

// Implementasi mock Publish
func (m *EventPublisherMock) Publish(routingKey string, data interface{}) error {
	args := m.Called(routingKey, data)
	return args.Error(0)
}
