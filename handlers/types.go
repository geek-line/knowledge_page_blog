package handlers

//DetailPage userのナレッジ詳細ページの要素
type DetailPage struct {
	Id               int
	Title            string
	Content          string
	SelectedTagNames []string
	UpdatedAt        string
	Likes            int
	EyeCatchSrc      string
}

//Header headerの要素
type Header struct {
	IsLogin bool
}

//Knowledges 各ナレッジの要素
type Knowledges struct {
	Id          int
	Title       string
	Content     string
	EyeCatchSrc string
}

//Tag 各タグの要素
type Tag struct {
	Id   int
	Name string
}

//IndexPage ナレッジ一覧ページの要素
type IndexPage struct {
	Id               int    //タイトル
	Title            string //タイトルの中身
	SelectedTagNames []string
	UpdatedAt        string
	Likes            int
	EyeCatchSrc      string
}
