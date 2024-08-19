package homeview

import (
	"os"
	"path"
	"strconv"

	"gioui.org/widget"
	"github.com/homearchbishop/downgram/internal/tg"
)

var (
	isSearching = false

	curInfo *tg.ChatInfo
)

func search(searchMode int, searchStr string) {
	if isSearching {
		return
	}
	if searchStr == "" {
		return
	}
	isSearching = true

	go func() {
		imgByteList.Overwrite(&[][]byte{})
		mediaInfoList.Overwrite(&[]tg.MediaInfo{})
		selectedList.Overwrite(&[]bool{})
		detailTipStr.Set("")

		searchTipStr.Set("searching...")
		info, err := getChatInfo(searchMode, searchStr)
		if err != nil {
			searchTipStr.Set(err.Error())
			isSearching = false
			return
		}
		curInfo = info
		searchTipStr.Set(
			"ID: " + strconv.FormatInt(info.ID, 10) + "\nType: " + info.Type,
		)

		detailTipStr.Set("updating media list...")
		if err := updateNextMediaList(info, imgByteList.Len(), 40); err != nil {
			detailTipStr.Set(err.Error())
			isSearching = false
			return
		}
		detailTipStr.Set("updated, total: " + strconv.Itoa(imgByteList.Len()))

		isSearching = false
	}()
}

func updateNextMediaListByBtn() {
	if curInfo == nil {
		return
	}
	if isSearching {
		return
	}
	isSearching = true

	go func() {
		detailTipStr.Set("updating media list...")
		if err := updateNextMediaList(curInfo, imgByteList.Len(), 40); err != nil {
			detailTipStr.Set(err.Error())
			isSearching = false
			return
		}
		detailTipStr.Set("updated, total: " + strconv.Itoa(imgByteList.Len()))

		isSearching = false
	}()
}

func downloadMediaToDir(dirPath string, mediaInfoList []tg.MediaInfo, selectedList []bool) {
	if isSearching {
		return
	}
	isSearching = true

	if dirPath == "" {
		isSearching = false
		return
	}

	if _, err := os.Stat(dirPath); os.IsNotExist(err) {
		os.MkdirAll(dirPath, os.ModePerm)
	}

	mediaInfoToDownload := []tg.MediaInfo{}
	for i, selected := range selectedList {
		if selected {
			mediaInfoToDownload = append(mediaInfoToDownload, mediaInfoList[i])
		}
	}

	go func() {
		detailTipStr.Set("downloading media, total: " + strconv.Itoa(len(mediaInfoToDownload)))
		done := make(chan error, len(mediaInfoToDownload))
		for i, media := range mediaInfoToDownload {
			go func(i int, media tg.MediaInfo) {
				// get size
				sizeType := ""
				if media.MediaType == "photo" {
					var maxSizeInfo *tg.PhotoSizeInfo
					for _, size := range media.PhotoSizes {
						if maxSizeInfo == nil || size.Chunks[len(size.Chunks)-1] > maxSizeInfo.Chunks[len(maxSizeInfo.Chunks)-1] {
							maxSizeInfo = &size
						}
					}
					sizeType = maxSizeInfo.SizeType
				}
				// get file
				filenamePrefix := strconv.FormatInt(media.MediaID, 10) + "_" + media.MediaType
				filenameSuffix := ""
				file, err := os.Create(path.Join(dirPath, filenamePrefix))
				if err != nil {
					done <- err
					return
				}
				// start download
				for offsetNum := 0; ; offsetNum++ {
					mediaByte, err := tg.GetMediaByteChunck(media, sizeType, offsetNum)
					if err != nil {
						done <- err
						file.Close()
						os.Remove(path.Join(dirPath, filenamePrefix))
						return
					}
					if len(mediaByte) == 0 {
						done <- nil
						break
					}
					if offsetNum == 0 {
						filenameSuffix = getFileSuffix(mediaByte)
					}
					file.Write(mediaByte)
				}
				// rename file
				file.Close()
				os.Rename(path.Join(dirPath, filenamePrefix), path.Join(dirPath, filenamePrefix+filenameSuffix))
			}(i, media)
		}
		successCount := 0
		var oneOfTheErr error
		for count := 0; count < len(mediaInfoToDownload); count++ {
			if err := <-done; err == nil {
				successCount++
			} else {
				oneOfTheErr = err
			}
		}
		detailTipStr.Set("downloaded media: " + strconv.Itoa(successCount) + ", failed: " +
			strconv.Itoa(len(mediaInfoToDownload)-successCount))
		if oneOfTheErr != nil {
			detailTipStr.Set(detailTipStr.Get() + oneOfTheErr.Error())
		}
		isSearching = false
	}()
}

func getChatInfo(searchMode int, searchStr string) (*tg.ChatInfo, error) {
	// get info
	var info *tg.ChatInfo
	var infoErr error
	if searchMode == idMode {
		id, _ := strconv.ParseInt(searchStr, 10, 64)
		info, infoErr = tg.GetChatInfoByID(id)
	} else {
		info, infoErr = tg.GetChatInfoByUsername(searchStr)
	}
	if infoErr != nil {
		return nil, infoErr
	}
	return info, nil
}

func getMediaHistoryWithOffsetLimit(info *tg.ChatInfo, offset, limit int) ([]tg.MediaInfo, error) {
	// get history
	history, err := tg.GetHistory(info.Type, info.ID, info.AccessHash, offset, limit)
	if err != nil {
		return nil, err
	}
	return history, nil
}

type getBytesChanSt struct {
	i     int
	bytes []byte
	err   error
}

func updateNextMediaList(info *tg.ChatInfo, offset, limit int) error {
	mediaList, err := getMediaHistoryWithOffsetLimit(info, offset, limit)
	if err != nil {
		return err
	}
	imgByteList.AppendSlice(make([][]byte, len(mediaList)))
	mediaInfoList.AppendSlice(mediaList)
	selectedList.AppendSlice(make([]bool, len(mediaList)))

	if len(mediaItemClickableList) < len(mediaList)+offset {
		for i := len(mediaItemClickableList); i < len(mediaList)+offset; i++ {
			mediaItemClickableList = append(mediaItemClickableList, &widget.Clickable{})
		}
	}

	done := make(chan getBytesChanSt, len(mediaList))
	for i, media := range mediaList {
		go func(i int, media tg.MediaInfo) {
			mediaByte, err := tg.GetFullMediaByte(media, "s")
			if err != nil {
				done <- getBytesChanSt{i, nil, err}
				return
			}
			done <- getBytesChanSt{i, mediaByte, nil}
		}(i, media)
	}
	for count := 0; count < len(mediaList); count++ {
		result := <-done
		if result.err != nil {
			continue
		}
		imgByteList.Set(result.i+offset, result.bytes)
	}
	close(done)
	return nil
}
