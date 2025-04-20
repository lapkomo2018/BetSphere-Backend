package service

import (
	"context"
	"errors"
	"time"

	"stavki/internal/database"
	"stavki/internal/model"

	"github.com/sirupsen/logrus"
)

type ChatService struct {
	msgDB       *database.MessageRepository
	userService *UserService
	hubs        map[uint64]*model.ChatHub
}

func NewChatService(msgDB *database.MessageRepository, userService *UserService) *ChatService {
	return &ChatService{
		msgDB:       msgDB,
		userService: userService,
		hubs:        make(map[uint64]*model.ChatHub),
	}
}

func (c *ChatService) HandleChatConnection(ctx context.Context, client *model.ChatClient, eventID uint64) error {
	if _, exists := c.hubs[eventID]; !exists {
		c.hubs[eventID] = model.NewChatHub(eventID)
	}

	user, err := c.userService.Get(ctx, client.UserID)
	if err != nil {
		client.Close()
		return err
	}

	hub := c.hubs[eventID]
	hub.AddClient(client)

	go func() {
		defer func() {
			client.Close()
			hub.RemoveClient(client)
		}()

		for {
			select {
			case <-client.Done():
				return
			case msg := <-client.MessageChannel():
				go c.HandleMessage(ctx, hub, user, msg)
			}
		}
	}()

	return nil
}

func (c *ChatService) HandleMessage(ctx context.Context, hub *model.ChatHub, user *model.User, msg model.ChatMessage) {
	// Fill in the message details
	msg.Username = user.Username
	msg.UserID = user.ID
	msg.Timestamp = time.Now()

	// Choose the appropriate function based on the action
	var fn func(ctx context.Context, hub *model.ChatHub, msg model.ChatMessage) error
	switch msg.Action {
	case model.ChatActionAddMessage:
		fn = c.AddMessage
	case model.ChatActionEditMessage:
		fn = c.EditMessage
	case model.ChatActionDeleteMessage:
		fn = c.DeleteMessage
	default:
		fn = func(ctx context.Context, hub *model.ChatHub, msg model.ChatMessage) error {
			return errors.New("unknown action")
		}
	}

	if err := fn(ctx, hub, msg); err != nil {
		logrus.WithFields(logrus.Fields{
			"userID":  user.ID,
			"eventID": hub.EventID(),
			"error":   err,
		}).Error("failed to handle message")
	}
}

func (c *ChatService) AddMessage(ctx context.Context, hub *model.ChatHub, msg model.ChatMessage) error {
	message, err := c.msgDB.Create(ctx, &model.Message{
		EventID:   hub.EventID(),
		UserID:    msg.UserID,
		Message:   msg.Message,
		CreatedAt: msg.Timestamp,
		UpdatedAt: msg.Timestamp,
	})
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"userID":  msg.UserID,
			"eventID": hub.EventID(),
			"error":   err,
		}).Error("failed to save message to database")
		return err
	}
	msg.MessageID = message.ID

	hub.SendMessage(msg)
	return nil
}

func (c *ChatService) EditMessage(ctx context.Context, hub *model.ChatHub, msg model.ChatMessage) error {
	message, err := c.msgDB.Get(ctx, msg.MessageID)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"userID":  msg.UserID,
			"eventID": hub.EventID(),
			"error":   err,
		}).Error("failed to get message from database")
		return err
	}

	if message.UserID != msg.UserID {
		logrus.WithFields(logrus.Fields{
			"userID":    msg.UserID,
			"eventID":   hub.EventID(),
			"messageID": msg.MessageID,
			"error":     errors.New("user not authorized to edit this message"),
		}).Error("unauthorized message edit attempt")
		return errors.New("unauthorized message edit attempt")
	}

	message.Message = msg.Message
	if _, err := c.msgDB.Save(ctx, message); err != nil {
		logrus.WithFields(logrus.Fields{
			"userID":  msg.UserID,
			"eventID": hub.EventID(),
			"error":   err,
		}).Error("failed to update message in database")
		return err
	}

	hub.SendMessage(msg)
	return nil
}

func (c *ChatService) DeleteMessage(ctx context.Context, hub *model.ChatHub, msg model.ChatMessage) error {
	message, err := c.msgDB.Get(ctx, msg.MessageID)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"userID":  msg.UserID,
			"eventID": hub.EventID(),
			"error":   err,
		}).Error("failed to get message from database")
		return err
	}

	if message.UserID != msg.UserID {
		logrus.WithFields(logrus.Fields{
			"userID":    msg.UserID,
			"eventID":   hub.EventID(),
			"messageID": msg.MessageID,
			"error":     errors.New("user not authorized to delete this message"),
		}).Error("unauthorized message delete attempt")
		return errors.New("unauthorized message delete attempt")
	}

	if err := c.msgDB.Delete(ctx, message); err != nil {
		logrus.WithFields(logrus.Fields{
			"userID":  msg.UserID,
			"eventID": hub.EventID(),
			"error":   err,
		}).Error("failed to delete message from database")
		return err
	}

	hub.SendMessage(msg)
	return nil
}
