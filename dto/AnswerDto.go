package dto

type AnswerShortDto struct {
	Id           int64           `json:"id"`
	Content      string          `json:"content"`
	ViewCount    int             `json:"view_count"`
	UpvoteCount  int             `json:"upvote_count"`
	CommentCount int             `json:"comment_count"`
	Creator      *UserProfileDto `json:"creator"`
}

type AnswerDetailDto struct {
	Id            int64           `json:"id"`
	Content       string          `json:"content"`
	ViewCount     int             `json:"view_count"`
	UpvoteCount   int             `json:"upvote_count"`
	CommentCount  int             `json:"comment_count"`
	CreateAt      int64           `json:"create_at"`
	UpdateAt      int64           `json:"update_at"`
	QuestionId    int64           `json:"question_id"`
	Creator       *UserProfileDto `json:"creator"`
}

type AnswerCreateDto struct {
	QuestionId int64 `json:"qid"`
	Content string `json:"content"`
}
