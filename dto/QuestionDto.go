package dto

type QuestionShortDto struct {
	Id          int64  `json:"id"`
	Title       string `json:"title"`
	Content     string `json:"content"`
	AnswerCount int    `json:"answer_count"`
	ViewCount   int    `json:"view_count"`
	UpdateAt    int64  `json:"update_at"`
}

type QuestionDetailDto struct {
	Id           int64           `json:"id"`
	Title        string          `json:"title"`
	Content      string          `json:"content"`
	Tag          string          `json:"tag"`
	AnswerCount  int             `json:"answer_count"`
	CommentCount int             `json:"comment_count"`
	ViewCount    int             `json:"view_count"`
	LikeCount    int             `json:"like_count"`
	CreateAt     int64           `json:"create_at"`
	UpdateAt     int64           `json:"update_at"`
	Creator      *UserProfileDto `json:"creator"`
}

type QuestionCreateDto struct {
	Title   string `json:"title"`
	Content string `json:"content"`
	Tag     string `json:"tag"`
}

type HotQuestionDto struct {
	Id          int64           `json:"id"`
	Title       string          `json:"title"`
	Content     string          `json:"content"`
	Heat        int             `json:"heat"`
	AnswerCount int             `json:"answer_count"`
	ViewCount   int             `json:"view_count"`
	UpdateAt    int64           `json:"update_at"`
	Answer      *AnswerShortDto `json:"answer"`
}
