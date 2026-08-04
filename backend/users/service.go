package users

import (
	"context"

	"github.com/clerk/clerk-sdk-go/v2"
	"github.com/clerk/clerk-sdk-go/v2/user"
)

func GetUserDetails(userID string) (*clerk.User, error) {
	clerk.SetKey("sk_test_XsiGhkWZeft81IydOVTHATS9UQZcrNDXckEPlym9M6")
	return user.Get(context.Background(), userID)
}
