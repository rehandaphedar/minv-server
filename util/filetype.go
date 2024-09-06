package util

func IsVideoFile(filetype string) bool {
	return filetype[:5] == "video"
}
