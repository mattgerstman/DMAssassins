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

// Marks a kill post as used for a game
func (game *Game) MarkPostUsed(post *KillPost) (appErr *ApplicationError) {
	res, err := db.Exec(`UPDATE dm_post_game_mapping SET used = true WHERE post_id = $1 AND game_id = $2`, post.PostId.String(), game.GameId.String())
	if err != nil {
		return NewApplicationError("Internal Error", err, ErrCodeDatabase)
	}

	// Check how many rows were affected by the update
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return NewApplicationError("Internal Error", err, ErrCodeDatabase)
	}

	// If the update succeeded no need to continue
	if rowsAffected != 0 {
		return nil
	}

	// If the update failed insert the post for the game
	res, err = db.Exec(`INSERT INTO dm_post_game_mapping (post_id, game_id, used) VALUES ($1, $2, true)`, post.PostId.String(), game.GameId.String())
	if err != nil {
		return NewApplicationError("Internal Error", err, ErrCodeDatabase)
	}
	return nil

}

// Get a random facebook post for a kill post
func (game *Game) GetRandomKillPost(assassin, target bool) (post *KillPost, appErr *ApplicationError) {
	var postIdBuffer, message string
	var official bool
	err := db.QueryRow(`SELECT post_id, message, official FROM dm_posts WHERE assassin = $1 AND target = $2 AND post_id NOT IN (SELECT post_id FROM dm_post_game_mapping WHERE game_id = $3 AND used = true OR allowed = false) ORDER BY RANDOM() LIMIT 1`, assassin, target, game.GameId.String()).Scan(&postIdBuffer, &message, &official)
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

	// If neither allows the post don't bother with any more queries
	if !allowAssassin && !allowTarget {
		return nil
	}

	post, appErr := game.GetRandomKillPost(allowAssassin, allowTarget)
	if appErr != nil {
		return appErr
	}

	message := post.Message

	message = strings.Replace(message, `ASSASSIN`, assassin.Username, -1)
	message = strings.Replace(message, `TARGET`, target.Username, -1)

	appErr = game.FacebookPost(message)
	if appErr != nil {
		return appErr
	}

	return game.MarkPostUsed(post)
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
