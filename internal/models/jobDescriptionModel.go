package models

type JobDescriptionTable struct {
	URL         string `json:"url"`
	ContentInfo string `json:"contentInfo"`
	TextContent string `json:"textContent"`
	UserName    string `json:"userName"`
}
