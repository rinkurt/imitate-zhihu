package dto

type QuestionShortDto struct {
	Id          int             `json:"id"`
	Title       string          `json:"title"`
	AnswerCount int             `json:"answer_count"`
	ViewCount   int             `json:"view_count"`
	GmtModified int64           `json:"gmt_modified"`
	BestAnswer  *AnswerShortDto `json:"best_answer"`
}

type QuestionDetailDto struct {
	Id           int      `json:"id"`
	Title        string   `json:"title"`
	Description  string   `json:"description"`
	Tag          string   `json:"tag"`
	AnswerCount  int      `json:"answer_count"`
	CommentCount int      `json:"comment_count"`
	ViewCount    int      `json:"view_count"`
	LikeCount    int      `json:"like_count"`
	GmtCreate    int64    `json:"gmt_create"`
	GmtModified  int64    `json:"gmt_modified"`
	Creator      *UserDto `json:"creator"`
}

type QuestionCreateDto struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Tag         string `json:"tag"`
}
