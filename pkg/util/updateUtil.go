package util

func GetStringOrDefault(input *string) string {
	if input != nil {
		return *input
	}
	return ""
}
