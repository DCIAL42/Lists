package follow

import "github.com/DCIAL42/lists/cmn"

type Follow struct {
	cmn.Model
	Follower string
	Followed string
}

type FollowResponse struct {
	ID       uint `json:"id"`
	Followed bool `json:"followed"`
}
