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
			var game_id = null;
			if (app.Session.get('game'))
			{
				game_id = game_id;
			}
			this.url = config.WEB_ROOT+'game/'+game_id+'/rules/'	
		}
		

	})
})();