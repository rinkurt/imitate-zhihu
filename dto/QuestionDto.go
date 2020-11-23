package dto

type QuestionShortDto struct {
	Id          int             `json:"id"`
	Title       string          `json:"title"`
	AnswerCount int             `json:"answer_count"`
	ViewCount   int             `json:"view_count"`
	UpdateAt    int64           `json:"update_at"`
	BestAnswer  *AnswerShortDto `json:"best_answer"`
}

type QuestionDetailDto struct {
	Id           int      `json:"id"`
	Title        string   `json:"title"`
	Content      string   `json:"content"`
	Tag          string   `json:"tag"`
	AnswerCount  int      `json:"answer_count"`
	CommentCount int      `json:"comment_count"`
	ViewCount    int      `json:"view_count"`
	LikeCount    int      `json:"like_count"`
	CreateAt     int64    `json:"create_at"`
	UpdateAt     int64    `json:"update_at"`
	Creator      *UserDto `json:"creator"`
}

type QuestionCreateDto struct {
	Title   string `json:"title"`
	Content string `json:"content"`
	Tag     string `json:"tag"`
}
