package models

type Metadata struct {
	FragmentsCount  int      `json:"fragments_count"`
	Categories      []string `json:"categories"`
	CategoriesCount int      `json:"categories_count"`
	Authors         []string `json:"authors"`
	AuthorsCount    int      `json:"authors_count"`
	Languages       []string `json:"languages"`
	LanguagesCount  int      `json:"languages_count"`
	ProjectInfo     struct {
		Source      string `json:"source"`
		ScrapeDate  string `json:"scrape_date"`
		APIVersion  string `json:"api_version"`
		Description string `json:"description"`
		License     string `json:"license"`
	} `json:"project_info"`
	HeteronymsInfo struct {
		MainHeteronyms []string `json:"main_heteronyms"`
		Description    string   `json:"description"`
	} `json:"heteronyms_info"`
}
