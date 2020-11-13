package repository

import (
	"imitate-zhihu/tool"
)

type Question struct {
	Id int `json:"id"`
	Title string `json:"title"`
}

func SelectQuestionById(id int) Question {
	db := tool.GetDatabase()
	question := Question{}
	db.First(&question, id)
	return question
}

func CreateQuestion(question Question) {
	db := tool.GetDatabase()
	db.Create(&question)
}
