package users

import "github.com/DCIAL42/lists/cmn"

type Following struct {
	cmn.Model
	Follower string `json:"from"`
	Followed string `json:"to"`
}

type UserResponse struct {
	ID       string `json:"id"`
	Username string `json:"username"`
}
