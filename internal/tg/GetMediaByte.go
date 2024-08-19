package tg

import (
	"github.com/gotd/td/tg"
)

func getMediaByteChunck(info MediaInfo, sizeType string, offset int64, limit int) ([]byte, error) {
	var location tg.InputFileLocationClass

	if info.MediaType == "photo" {
		location = &tg.InputPhotoFileLocation{
			ID:            info.MediaID,
			AccessHash:    info.AccessHash,
			ThumbSize:     sizeType,
			FileReference: info.FileReference,
		}
	} else if info.MediaType == "document" {
		location = &tg.InputDocumentFileLocation{
			ID:            info.MediaID,
			AccessHash:    info.AccessHash,
			ThumbSize:     sizeType,
			FileReference: info.FileReference,
		}
	}

	fileResp, err := client.API().UploadGetFile(*clientCtx, &tg.UploadGetFileRequest{
		Location: location,
		Offset:   offset,
		Limit:    limit,
	})
	if err != nil {
		return nil, err
	}

	return fileResp.(*tg.UploadFile).Bytes, nil
}

func GetFullMediaByte(info MediaInfo, sizeType string) ([]byte, error) {
	chunckLimit := 4 * 1024 // 4kB
	chunckOffset := int64(0)
	mediaByte := []byte{}
	for {
		chunckByte, err := getMediaByteChunck(info, sizeType, chunckOffset, chunckLimit)
		if len(chunckByte) == 0 {
			break
		}
		if err != nil {
			return nil, err
		}
		mediaByte = append(mediaByte, chunckByte...)
		chunckOffset += int64(chunckLimit)
	}
	return mediaByte, nil
}

func GetMediaByteChunck(info MediaInfo, sizeType string, offsetNum int) ([]byte, error) {
	chunckLimit := 512 * 1024 // 512kB
	chunckOffset := int64(chunckLimit * offsetNum)
	chunckByte, err := getMediaByteChunck(info, sizeType, chunckOffset, chunckLimit)
	if err != nil {
		return nil, err
	}
	return chunckByte, nil
}
