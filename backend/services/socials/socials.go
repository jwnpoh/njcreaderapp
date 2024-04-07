package socials

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/jwnpoh/njcreaderapp/backend/internal/core"
	"github.com/jwnpoh/njcreaderapp/backend/services/serializer"
	"github.com/jwnpoh/njcreaderapp/backend/services/users"
)

type SocialsDB interface {
	GetFollowing(userID uuid.UUID) ([]uuid.UUID, error)
	GetFollowedBy(userID uuid.UUID) ([]uuid.UUID, error)
	Follow(userID, toFollow uuid.UUID) error
	UnFollow(userID, toUnFollow uuid.UUID) error
	Like(userID, postID uuid.UUID) error
	Unlike(userID, postID uuid.UUID) error
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

func (sDB *Socials) GetFriends(userID uuid.UUID) (serializer.Serializer, error) {
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

func (sDB *Socials) Follow(userID, toFollow uuid.UUID, follow bool) (serializer.Serializer, error) {
	if !follow {
		return sDB.unfollow(userID, toFollow)
	}

	toFollowUser, err := sDB.GetUser("id", toFollow)
	if err != nil {
		return nil, fmt.Errorf("unable to get user info to follow - %w", err)
	}

	err = sDB.db.Follow(userID, toFollow)
	if err != nil {
		return nil, fmt.Errorf("unable to follow user %s - %w", toFollowUser.DisplayName, err)
	}

	return serializer.NewSerializer(false, fmt.Sprintf("successfully followed user %s", toFollowUser.DisplayName), nil), nil
}

func (sDB *Socials) unfollow(userID, toUnFollow uuid.UUID) (serializer.Serializer, error) {
	toUnFollowUser, err := sDB.GetUser("id", toUnFollow)
	if err != nil {
		return nil, fmt.Errorf("unable to get user info to follow - %w", err)
	}

	err = sDB.db.UnFollow(userID, toUnFollow)
	if err != nil {
		return nil, fmt.Errorf("unable to follow user %s - %w", toUnFollowUser.DisplayName, err)
	}

	return serializer.NewSerializer(false, fmt.Sprintf("successfully unfollowed user %s", toUnFollowUser.DisplayName), nil), nil
}

func (sDB *Socials) Like(userID, postID uuid.UUID, like bool) (serializer.Serializer, error) {
	if !like {
		return sDB.unlike(userID, postID)
	}

	err := sDB.db.Like(userID, postID)
	if err != nil {
		return nil, fmt.Errorf("error liking post - %w", err)
	}

	return serializer.NewSerializer(false, "liked post", nil), nil
}

func (sDB *Socials) unlike(userID, postID uuid.UUID) (serializer.Serializer, error) {
	err := sDB.db.Unlike(userID, postID)
	if err != nil {
		return nil, fmt.Errorf("error unliking post - %w", err)
	}

	return serializer.NewSerializer(false, "unliked post", nil), nil
}
