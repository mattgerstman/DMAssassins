// Leaderboard model, displays high scores for a game
// js/models/leaderboard.js

var app = app || {Models:{}, Views:{}, Routers:{}, Running:{}, Session:{}};
(function() {
	'use strict';
	
	app.Models.Leaderboard = Backbone.Model.extend({
		defaults: {
			teams_enabled: true,
			col_width:33.3333,
			users : [
				{
					name:"Matt",
					kills:75,
					team_name:"Tech"
				},
				{
					name:"Jimmy",
					kills:5,
					team_name:"Morale"					
				}
			],
			teams: [
				{
					Tech: 4,
					Morale: 1
				}
			]
		},
		parse: function(response){
			var data = response.response;
			data.col_width = (100 / (2 + data.teams_enabled))
			return data;
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