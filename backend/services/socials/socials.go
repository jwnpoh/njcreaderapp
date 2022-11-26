package socials

import (
	"fmt"

	"github.com/jwnpoh/njcreaderapp/backend/internal/core"
	"github.com/jwnpoh/njcreaderapp/backend/services/serializer"
	"github.com/jwnpoh/njcreaderapp/backend/services/users"
)

type SocialsDB interface {
	GetFollowing(userID int) ([]int, error)
	GetFollowedBy(userID int) ([]int, error)
}

type Socials struct {
	db SocialsDB
	users.UsersDB
}

func NewSocialsDB(socialsDB SocialsDB, usersDB users.UsersDB) *Socials {
	return &Socials{
		db:      socialsDB,
		UsersDB: usersDB,
	}
}

func (sDB *Socials) GetFriends(userID int) (serializer.Serializer, error) {

	following, err := sDB.db.GetFollowing(userID)
	if err != nil {
		return nil, fmt.Errorf("unable to get following ids - %w", err)
	}

	followingUsers := make([]string, 0, len(following))
	for _, v := range following {
		user, err := sDB.GetUser("id", v)
		if err != nil {
			return nil, fmt.Errorf("unable to get following users - %w", err)
		}
		followingUsers = append(followingUsers, user.DisplayName)
	}

	followedBy, err := sDB.db.GetFollowedBy(userID)
	if err != nil {
		return nil, fmt.Errorf("unable to get followed by ids - %w", err)
	}

	followedByUsers := make([]string, 0, len(followedBy))
	for _, w := range following {
		user, err := sDB.GetUser("id", w)
		if err != nil {
			return nil, fmt.Errorf("unable to get followed by users - %w", err)
		}
		followedByUsers = append(followedByUsers, user.DisplayName)
	}

	friends := core.Relations{
		UserID:          userID,
		FollowingIDs:    following,
		FollowingUsers:  followingUsers,
		FollowedByIDs:   followedBy,
		FollowedByUsers: followedByUsers,
	}

	return serializer.NewSerializer(false, "successfully got friends", friends), nil
}
