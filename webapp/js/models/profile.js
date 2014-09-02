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
			this.idAttribute = 'username'
			this.url = config.WEB_ROOT  + 'users/' + this.get('user_id') + '/';
		}
	})
})();