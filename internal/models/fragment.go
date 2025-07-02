package models

type Fragment struct {
	ID      int    `json:"id"`
	URL     string `json:"url"`
	Title   string `json:"title"`
	Text    string `json:"text"`
	Length  int    `json:"length"`
	Excerpt string `json:"excerpt"`
}

type FragmentResponse struct {
	ID      int    `json:"id"`
	URL     string `json:"url"`
	Title   string `json:"title"`
	Text    string `json:"text"`
	Length  int    `json:"length"`
	Excerpt string `json:"excerpt"`
}

func (f *Fragment) ToResponse() FragmentResponse {
	return FragmentResponse{
		ID:      f.ID,
		URL:     f.URL,
		Title:   f.Title,
		Text:    f.Text,
		Length:  f.Length,
		Excerpt: f.Excerpt,
	}
}
