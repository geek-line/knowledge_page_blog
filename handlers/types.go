package handlers

//DetailPage userのナレッジ詳細ページの要素
type DetailPage struct {
	ID           int
	Title        string
	Content      string
	SelectedTags []Tag
	UpdatedAt    string
	Likes        int
	EyeCatchSrc  string
}

//Header headerの要素
type Header struct {
	IsLogin bool
}

//Knowledges 各ナレッジの要素
type Knowledges struct {
	ID          int
	Title       string
	Content     string
	EyeCatchSrc string
}

//Tag 各タグの要素
type Tag struct {
	ID   int
	Name string
}

//IndexElem ナレッジ一覧ページの各リストの要素
type IndexElem struct {
	ID           int    //タイトル
	Title        string //タイトルの中身
	SelectedTags []Tag
	CreatedAt    string
	UpdatedAt    string
	Likes        int
	EyeCatchSrc  string
}

//IndexPage ナレッジ一覧ページの全体の要素
type IndexPage struct {
	PageNationElems []float64
	IndexElems      []IndexElem
}
