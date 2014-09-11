// Games Model, Handles game creation, selection, and joining
// js/models/games.js
var app = app || {
	Models: {},
	Views: {},
	Routers: {},
	Running: {},
	Session: {}
};
(function() {
	'use strict';
	app.Models.SelectGames = Backbone.Model.extend({
	
		// default properties with a fake game
		defaults: {
			games: [{
				game_id: '',
				game_name: 'Triwizard Tournament',
				game_started: false,
				game_has_password: false
			}]
		},
		
		// handle on initiliazation
		initialize: function() {
			this.url = config.WEB_ROOT + 'game/'
		},
		
		// automatically called by fetch() to place a wrapper around the response
		parse: function(response) {
			var wrapper = {
				games: response.response
			};
			return wrapper;
		},
		
		// creates a new game, posts it to the server, and sets the current user's game to that game
		create: function(name, password) {
			var that = this;
			var data = {
				user_id: app.Session.get('user_id'),
				game_name: name,
				game_password: password
			}
			$.post(this.url, data, function(response) {
				app.Session.set('game', JSON.stringify(response.response))
				that.trigger('finish_set_game');
			})
		},
		
		// allows a user to join a given game
		join: function(game_id, user_id, password) {
			var that = this;
			var data = {
				game_password: password
			}
			var join_url = config.WEB_ROOT + 'game/' + game_id + '/users/' + user_id + '/';
			$.post(join_url, data, function(response) {
				app.Session.set('game', JSON.stringify(response.response))
				that.trigger('finish_set_game');
			})
		},
		setUser: function(user_id) {
			this.url = config.WEB_ROOT + 'users/' + user_id + '/game/new/';
		}
	})
})();