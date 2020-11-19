package dto

type AnswerShortDto struct {
	Id           int      `json:"id"`
	Answer       string   `json:"answer"`
	ViewCount    int      `json:"view_count"`
	UpvoteCount  int      `json:"upvote_count"`
	CommentCount int      `json:"comment_count"`
	Creator      *UserDto `json:"creator"`
}

type AnswerDetailDto struct {
	Id            int      `json:"id"`
	Answer        string   `json:"answer"`
	ViewCount     int      `json:"view_count"`
	UpvoteCount   int      `json:"upvote_count"`
	DownvoteCount int      `json:"downvote_count"`
	CommentCount  int      `json:"comment_count"`
	GmtCreate     int64    `json:"gmt_create"`
	GmtModified   int64    `json:"gmt_modified"`
	Creator       *UserDto `json:"creator"`
}
