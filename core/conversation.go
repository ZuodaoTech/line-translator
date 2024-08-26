package core

import (
	"context"
	"time"
)

const (
	ConversationStrategyZhEn = "zh,en"
	ConversationStrategyZhJa = "zh,ja"
)

type (
	Conversation struct {
		ID        uint64 `json:"id"`
		Channel   string `json:"channel"`  // line, slack, etc
		ChannelID string `json:"group_id"` // for line, it's group_id
		Strategy  string `json:"strategy"` // zh,en; zh,ja

		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
	}

	ConversationStore interface {
		// INSERT INTO @@table (
		//  channel, channel_id, strategy,
		//  created_at, updated_at
		// ) VALUES (
		//  @item.Channel, @item.ChannelID, @item.Strategy,
		//  NOW(), NOW()
		// ) ON CONFLICT (channel, channel_id) DO UPDATE SET
		//  strategy = @item.Strategy,
		//  updated_at = NOW()
		// RETURNING id;
		CreateConversation(ctx context.Context, item *Conversation) (uint64, error)

		// SELECT * FROM @@table
		// WHERE channel = @channel AND channel_id = @channelID
		// LIMIT 1;
		GetConversationByChannel(ctx context.Context, channel, channelID string) (*Conversation, error)

		// UPDATE @@table
		// SET
		//  strategy = @strategy,
		//  updated_at = NOW()
		// WHERE id = @id;
		UpdateConversationStrategy(ctx context.Context, id uint64, strategy string) error
	}
)
