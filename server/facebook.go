package main

import (
	fb "github.com/huandu/facebook"
)

func facebook() (interface{}, *ApplicationError) {
	
	username := "Matt"
	var facebook_id, facebook_token string;
	_ = db.QueryRow(`SELECT facebook_id, facebook_token FROM dm_users WHERE username = $1`, username).Scan(&facebook_id, &facebook_token)


	fb.Version = "v2.0"
	var app = fb.New("643600385736129", "73cbc95ae6de7a6c26b16318330f796a")
	app.RedirectUri = "http://dmassassins.com"

	session := app.Session(facebook_token)

	//path := "/" + facebook_id + "/friends/"
	
	res, _ := session.Get("/me/friends/", fb.Params{

	})
	
	return res, nil
}
