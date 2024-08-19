package homeview

var magicNumbers = map[string]struct {
	suffix string
	offset int
}{
	"\xFF\xD8\xFF":     {suffix: ".jpg", offset: 0},
	"\x89PNG":          {suffix: ".png", offset: 0},
	"GIF87a":           {suffix: ".gif", offset: 0},
	"GIF89a":           {suffix: ".gif", offset: 0},
	"WEBP":             {suffix: ".webp", offset: 8},
	"ftypmp42":         {suffix: ".mp4", offset: 4},
	"ftypisom":         {suffix: ".mp4", offset: 4},
	"\x1A\x45\xDF\xA3": {suffix: ".mkv", offset: 0},
	"ftypqt":           {suffix: ".mov", offset: 4},
}

func getFileSuffix(buf []byte) string {
	for magic, ext := range magicNumbers {
		if len(buf) >= ext.offset+len(magic) && string(buf[ext.offset:ext.offset+len(magic)]) == magic {
			return ext.suffix
		}
	}

	return ""
}
