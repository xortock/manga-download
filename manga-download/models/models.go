package models

type Chapter struct {
	Number float64
	Id int
}

type ChapterApiResult struct {
	Status int `json:"status"`
	Result ChapterApiImageResult `json:"result"`
}

type ChapterApiImageResult struct {
	Images []interface{} `json:"images"`
}

type ChaptersApiResult struct {
	Status int `json:"status"`
	HtmlResult ChaptersApiHtmlResult `json:"result"`

}

type ChaptersApiHtmlResult struct {
	HtmlContent string `json:"html"`
    TitleFormat string `json:"title_format"`
}