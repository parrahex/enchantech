package app

type pageData struct {
	Profile profile
}

func newPageData(content profile) pageData {
	return pageData{Profile: content}
}
