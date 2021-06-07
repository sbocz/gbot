package util

import "github.com/diamondburned/arikawa/v2/discord"

func UserIdFromBytes(bytes []byte) (discord.UserID, error) {
	userSf, err := discord.ParseSnowflake(string(bytes))
	if err != nil {
		return 0, err
	}
	return discord.UserID(userSf), nil
}

func UserIdsFromBytes(bytesSlice [][]byte) ([]discord.UserID, error) {
	userIds := make([]discord.UserID, 0)

	for _, bytes := range bytesSlice {
		userSf, err := UserIdFromBytes(bytes)
		if err != nil {
			return nil, err
		}
		userIds = append(userIds, discord.UserID(userSf))
	}
	return userIds, nil
}
