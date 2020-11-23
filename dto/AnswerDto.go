package dto

type AnswerShortDto struct {
	Id           int      `json:"id"`
	Content      string   `json:"content"`
	ViewCount    int      `json:"view_count"`
	UpvoteCount  int      `json:"upvote_count"`
	CommentCount int      `json:"comment_count"`
	Creator      *UserDto `json:"creator"`
}

type AnswerDetailDto struct {
	Id            int      `json:"id"`
	Content       string   `json:"content"`
	ViewCount     int      `json:"view_count"`
	UpvoteCount   int      `json:"upvote_count"`
	DownvoteCount int      `json:"downvote_count"`
	CommentCount  int      `json:"comment_count"`
	CreateAt      int64    `json:"create_at"`
	UpdateAt      int64    `json:"update_at"`
	Creator       *UserDto `json:"creator"`
}
