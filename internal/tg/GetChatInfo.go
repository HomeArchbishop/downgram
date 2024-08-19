package tg

import (
	"errors"

	"github.com/gotd/td/tg"
)

type ChatInfo struct {
	ID         int64
	AccessHash int64
	Type       string // "channel" or "chat" or "user"
}

func GetChatInfoByID(id int64) (*ChatInfo, error) {
	// group & channel will be treated as channel
	channelResp, err := client.API().ChannelsGetChannels(*clientCtx, []tg.InputChannelClass{&tg.InputChannel{
		ChannelID:  id,
		AccessHash: 0,
	}})
	if err != nil {
		return nil, err
	}
	chat := channelResp.(*tg.MessagesChats).Chats
	if len(chat) == 0 {
		return nil, errors.New("channel not found")
	}
	return &ChatInfo{
		ID:         chat[0].(*tg.Channel).ID,
		AccessHash: chat[0].(*tg.Channel).AccessHash,
		Type:       "channel",
	}, nil
}

func GetChatInfoByUsername(username string) (*ChatInfo, error) {
	// channel & chat will be treated as channel
	resp, err := client.API().ContactsResolveUsername(*clientCtx, username)
	if err != nil {
		return nil, err
	}
	switch resp.Peer.(type) {
	case *tg.PeerChannel:
		return &ChatInfo{
			ID:         resp.Peer.(*tg.PeerChannel).ChannelID,
			AccessHash: resp.Chats[0].(*tg.Channel).AccessHash,
			Type:       "channel",
		}, nil
	case *tg.PeerUser:
		return &ChatInfo{
			ID:         resp.Peer.(*tg.PeerUser).UserID,
			AccessHash: resp.Users[0].(*tg.User).AccessHash,
			Type:       "user",
		}, nil
	}

	return nil, errors.New("Unknown chat type")
}
