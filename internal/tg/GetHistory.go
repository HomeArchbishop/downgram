package tg

import (
	"log"

	"github.com/gotd/td/tg"
)

type MediaInfo struct {
	MediaType     string // photo, document
	MediaID       int64
	MessageID     int
	AccessHash    int64
	FileReference []byte
	PhotoSizes    []PhotoSizeInfo
	VideoSize     int64
	Date          int
	Duration      float64
}

type PhotoSizeInfo struct {
	SizeType string // s, m, x, y, z, a, b, c, d
	Width    int
	Height   int
	Chunks   []int // increase step by step
}

func GetHistory(chatType string, id, accessHash int64, offset, limit int) ([]MediaInfo, error) {
	var historyReqPeer tg.InputPeerClass
	var historyMessages []tg.MessageClass
	if chatType == "channel" {
		historyReqPeer = &tg.InputPeerChannel{
			ChannelID:  id,
			AccessHash: accessHash,
		}
	} else if chatType == "user" {
		historyReqPeer = &tg.InputPeerUser{
			UserID:     id,
			AccessHash: accessHash,
		}
	}
	historyReq := &tg.MessagesSearchRequest{
		Peer:      historyReqPeer,
		Filter:    &tg.InputMessagesFilterPhotoVideo{},
		Limit:     limit,
		AddOffset: offset,
	}
	historyResp, err := client.API().MessagesSearch(*clientCtx, historyReq)
	if err != nil {
		log.Println("Error:", err)
		return nil, err
	}
	if chatType == "channel" {
		historyMessages = historyResp.(*tg.MessagesChannelMessages).Messages
	} else {
		historyMessages = historyResp.(*tg.MessagesMessagesSlice).Messages
	}

	mediaInfoList := []MediaInfo{}

	for _, message := range historyMessages {
		if msg, ok := message.(*tg.Message); ok {
			switch msg.Media.(type) {
			case *tg.MessageMediaPhoto:
				photoSizeInfoList := []PhotoSizeInfo{}
				for _, photoSize := range msg.Media.(*tg.MessageMediaPhoto).Photo.(*tg.Photo).Sizes {
					switch photoSize.(type) {
					case *tg.PhotoSize:
						photoSizeInfoList = append(photoSizeInfoList, PhotoSizeInfo{
							SizeType: photoSize.(*tg.PhotoSize).Type,
							Width:    photoSize.(*tg.PhotoSize).W,
							Height:   photoSize.(*tg.PhotoSize).H,
							Chunks:   []int{photoSize.(*tg.PhotoSize).Size},
						})
					case *tg.PhotoSizeProgressive:
						photoSizeInfoList = append(photoSizeInfoList, PhotoSizeInfo{
							SizeType: photoSize.(*tg.PhotoSizeProgressive).Type,
							Width:    photoSize.(*tg.PhotoSizeProgressive).W,
							Height:   photoSize.(*tg.PhotoSizeProgressive).H,
							Chunks:   photoSize.(*tg.PhotoSizeProgressive).Sizes,
						})
					}
				}
				mediaInfoList = append(mediaInfoList, MediaInfo{
					MediaType:     "photo",
					MediaID:       msg.Media.(*tg.MessageMediaPhoto).Photo.(*tg.Photo).ID,
					MessageID:     msg.ID,
					AccessHash:    msg.Media.(*tg.MessageMediaPhoto).Photo.(*tg.Photo).AccessHash,
					FileReference: msg.Media.(*tg.MessageMediaPhoto).Photo.(*tg.Photo).FileReference,
					Date:          msg.Date,
					Duration:      0,
					PhotoSizes:    photoSizeInfoList,
				})
			case *tg.MessageMediaDocument:
				duration := float64(0)
				for _, attr := range msg.Media.(*tg.MessageMediaDocument).Document.(*tg.Document).Attributes {
					if attr, ok := attr.(*tg.DocumentAttributeVideo); ok {
						duration = attr.Duration
					}
				}
				photoSizeInfoList := []PhotoSizeInfo{}
				for _, photoSize := range msg.Media.(*tg.MessageMediaDocument).Document.(*tg.Document).Thumbs {
					switch photoSize.(type) {
					case *tg.PhotoSize:
						photoSizeInfoList = append(photoSizeInfoList, PhotoSizeInfo{
							SizeType: photoSize.(*tg.PhotoSize).Type,
							Width:    photoSize.(*tg.PhotoSize).W,
							Height:   photoSize.(*tg.PhotoSize).H,
							Chunks:   []int{photoSize.(*tg.PhotoSize).Size},
						})
					case *tg.PhotoSizeProgressive:
						photoSizeInfoList = append(photoSizeInfoList, PhotoSizeInfo{
							SizeType: photoSize.(*tg.PhotoSizeProgressive).Type,
							Width:    photoSize.(*tg.PhotoSizeProgressive).W,
							Height:   photoSize.(*tg.PhotoSizeProgressive).H,
							Chunks:   photoSize.(*tg.PhotoSizeProgressive).Sizes,
						})
					}
				}
				mediaInfoList = append(mediaInfoList, MediaInfo{
					MediaType:     "document",
					MediaID:       msg.Media.(*tg.MessageMediaDocument).Document.(*tg.Document).ID,
					MessageID:     msg.ID,
					AccessHash:    msg.Media.(*tg.MessageMediaDocument).Document.(*tg.Document).AccessHash,
					FileReference: msg.Media.(*tg.MessageMediaDocument).Document.(*tg.Document).FileReference,
					Date:          msg.Date,
					Duration:      duration,
					VideoSize:     msg.Media.(*tg.MessageMediaDocument).Document.(*tg.Document).Size,
					PhotoSizes:    photoSizeInfoList,
				})
			}
		}
	}
	return mediaInfoList, nil
}
