package handlers

//Tag 各タグの要素
type Tag struct {
	ID   int
	Name string
}

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

//IndexElem ナレッジ一覧ページの要素
type IndexElem struct {
	ID           int    //タイトル
	Title        string //タイトルの中身
	SelectedTags []Tag
	CreatedAt    string
	UpdatedAt    string
	Likes        int
	EyeCatchSrc  string
}

//Page ページネーションの際に用いるページの要素
type Page struct {
	PageNum  int
	IsSelect bool
}

//PageNation ページネーションの全体の要素
type PageNation struct {
	PageElems   []Page
	PageNum     int
	NextPageNum int
	PrevPageNum int
}

//IndexPage ナレッジ一覧ページの全体の要素
type IndexPage struct {
	PageNation PageNation
	IndexElems []IndexElem
}

//EyeCatch アイキャッチの要素
type EyeCatch struct {
	ID   int
	Name string
	Src  string
}
