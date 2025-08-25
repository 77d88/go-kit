package xpg

func WarpLike(str string) string {
	if str != "" {
		return "%" + str + "%"
	}
	return ""
}

func WarpLikeRight(str string) string {
	if str != "" {
		return str + "%"
	}
	return ""
}
func WarpLikeLeft(str string) string {
	if str != "" {
		return "%" + str
	}
	return ""
}
