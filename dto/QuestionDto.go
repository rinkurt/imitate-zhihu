package dto

type QuestionDto struct {
	Id           int      `json:"id"`
	Title        string   `json:"title"`
	Description  string   `json:"description"`
	Tag          string   `json:"tag"`
	CommentCount int      `json:"comment_count"`
	ViewCount    int      `json:"view_count"`
	LikeCount    int      `json:"like_count"`
	GmtCreate    int64    `json:"gmt_create"`
	GmtModified  int64    `json:"gmt_modified"`
	Creator      *UserDto `json:"creator"`
}
