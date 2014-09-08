// Model for user profile
// js/models/profile.js

var app = app || {Models:{}, Views:{}, Routers:{}, Running:{}, Session:{}};
(function() {
	'use strict';
	
	app.Models.Profile = Backbone.Model.extend({

		// default profile properties
		defaults: {
			'user_id' : '',
			'username' : '',
			'email' : '',
			'properties': {
				'name' : '',
				'facebook': '',
				'twitter': '',
				'photo_thumb' : SPY,
				'photo' : SPY
			}

		},
				
		// parses the api wrapper, automatically called by fetch()
		parse: function(response) {                
               return response.response;
        },
        
        // loaded on initialization
		initialize: function() {
			this.idAttribute = 'user_id';
			var game_id = app.Session.getGameId();
			this.url = config.WEB_ROOT + 'game/' + game_id + '/users/' + this.get('user_id') + '/';
		},
		changeGame: function() {
			var game_id = app.Session.getGameId();
			this.url = config.WEB_ROOT + 'game/' + game_id + '/users/' + this.get('user_id') + '/';			
		}
	})
})();