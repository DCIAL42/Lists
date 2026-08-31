package users

import (
	"context"
	"fmt"
	"os"

	"github.com/clerk/clerk-sdk-go/v2"
	"github.com/clerk/clerk-sdk-go/v2/user"
)

func GetUserDetails(userID string) (*clerk.User, error) {
	clerkSecret, ok := os.LookupEnv("CLERK_SECRET_KEY")

	if !ok {
		panic("clerk secret key not found")
	}

	clerk.SetKey(clerkSecret)

	return user.Get(context.Background(), userID)
}

func GetUserByUsername(username string) (*clerk.User, error) {
	clerkSecret, ok := os.LookupEnv("CLERK_SECRET_KEY")

	if !ok {
		panic("clerk secret key not found")
	}

	clerk.SetKey(clerkSecret)

	userList, err := user.List(context.Background(), &user.ListParams{
		UsernameQuery: &username,
	})

	if err != nil {
		return nil, err
	}

	if len(userList.Users) == 1 {
		return userList.Users[0], nil
	}

	return nil, fmt.Errorf("user not found")
}

func (s *Service) getAllUsers() ([]UserResponse, error) {
	clerkSecret, ok := os.LookupEnv("CLERK_SECRET_KEY")

	if !ok {
		panic("clerk secret key not found")
	}

	clerk.SetKey(clerkSecret)

	userList, err := user.List(context.Background(), &user.ListParams{})
	res := make([]UserResponse, 0, userList.TotalCount)
	for _, u := range userList.Users {
		res = append(res, UserResponse{u.ID, *u.Username})
	}
	return res, err
}
