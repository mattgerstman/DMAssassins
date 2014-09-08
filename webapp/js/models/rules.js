// Rules model, loads rules from the db so that admins can define custom rules per game
// js/models/user.js

var app = app || {Models:{}, Views:{}, Routers:{}, Running:{}, Session:{}};
(function() {
	'use strict';
	
	app.Models.Rules = Backbone.Model.extend({
		defaults: {
			rules : 'yo'
		},
		parse: function(response) {
			return { rules: response.response }
		},
		initialize: function(){
			var game_id = app.Session.getGameId();
			this.url = config.WEB_ROOT+'game/'+game_id+'/rules/'	
		},
		loadGame: function(){
			var game_id = app.Session.getGameId();
			this.url = config.WEB_ROOT + 'game/' + game_id + '/rules/'
			
		}
		

	})
})();