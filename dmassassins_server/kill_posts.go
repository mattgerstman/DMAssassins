package main

import (
	"code.google.com/p/go-uuid/uuid"
	"strings"
)

type KillPost struct {
	PostId   uuid.UUID `json:"post_id"`
	Message  string    `json:"message"`
	Official bool      `json:"official"`
	Assassin bool      `json:"assassin"`
	Target   bool      `json:"target"`
}

// Get a random facebook post for a kill post
func (game *Game) GetRandomKillPost(assassin, target bool) (post *KillPost, appErr *ApplicationError) {
	var postIdBuffer, message string
	var official bool
	err := db.QueryRow(`SELECT post_id, message, official FROM dm_posts WHERE assassin = $1 AND target = $2 AND post_id NOT IN (SELECT post_id FROM dm_post_game_mapping WHERE game_id = $3 AND used = false AND allowed = true) ORDER BY RANDOM() LIMIT 1`, assassin, target, game.GameId.String()).Scan(&postIdBuffer, &message, &official)
	if err != nil {
		return nil, NewApplicationError(`Internal Error`, err, ErrCodeDatabase)
	}
	postId := uuid.Parse(postIdBuffer)
	return &KillPost{postId, message, official, assassin, target}, nil
}

// post a kill from an assassin target pair
func (game *Game) PostKill(assassin, target *User) (appErr *ApplicationError) {
	allowAssassin, appErr := assassin.GetUserPropertyBool(`allow_post`)
	if appErr != nil {
		return appErr
	}
	allowTarget, appErr := target.GetUserPropertyBool(`allow_post`)
	if appErr != nil {
		return appErr
	}
	_ = allowAssassin
	_ = allowTarget

	post, appErr := game.GetRandomKillPost(true, true)
	if appErr != nil {
		return appErr
	}

	message := post.Message

	message = strings.Replace(message, `ASSASSIN`, assassin.Username, -1)
	message = strings.Replace(message, `TARGET`, target.Username, -1)

	return game.FacebookPost(message)
}

// Handles a kill post and gets the assassin/target/game structs to pass to it
func (assassin *User) HandleKillPost(gameId, oldTargetId uuid.UUID) (appErr *ApplicationError) {
	game, appErr := GetGameById(gameId)
	if appErr != nil {
		return appErr
	}
	target, appErr := GetUserById(oldTargetId)
	if appErr != nil {
		return appErr
	}
	return game.PostKill(assassin, target)
}
