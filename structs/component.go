package structs

//UserDetailPage userのナレッジ詳細ページの要素
type UserDetailPage struct {
	Knowledge    Knowledge
	SelectedTags []Tag
}

//Header headerの要素
type Header struct {
	IsLogin bool
}

//IndexElem ナレッジ一覧ページの要素
type IndexElem struct {
	Knowledge    Knowledge
	SelectedTags []Tag
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

//UserIndexPage ナレッジ一覧ページの全体の要素
type UserIndexPage struct {
	PageNation  PageNation
	IndexElems  []IndexElem
	CurrentSort string
	TagRanking  []TagRankingElem
}

//TagRankingElem ナレッジ一覧ページのタグランキングの各要素
type TagRankingElem struct {
	TagID              int
	TagName            string
	CountOfRefferenced int
}
