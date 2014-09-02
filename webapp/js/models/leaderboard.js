// Leaderboard model, displays high scores for a game
// js/models/leaderboard.js

var app = app || {Models:{}, Views:{}, Routers:{}, Running:{}, Session:{}};
(function() {
	'use strict';
	
	app.Models.Leaderboard = Backbone.Model.extend({
		defaults: {
			game_id : null,
			users : [
				{
					name:"Matt",
					kills:75
				},
				{
					name:"Jimmy",
					kills:5
				}
			],
			teams: [
				{
					name: "Slytherin",
					alive: 7,
					kills: 75
				},				
				{
					name: "Ravenclaw",
					alive: 7,
					kills: 14
				}
			]
		},
		parse: function(response){
			return {users: response.response};
		},
		initialize: function(){
			var game_id = null;
			if (app.Session.get('game'))
			{
				game_id = app.Session.get('game').game_id
			}
			this.url = config.WEB_ROOT + 'game/' + game_id + '/leaderboard/'
		},
		loadGame: function(){
			var game_id = null;
			if (app.Session.get('game'))
			{
				game_id = app.Session.get('game').game_id
			}
			this.url = config.WEB_ROOT + 'game/' + game_id + '/leaderboard/'
			
		}

	})
})();