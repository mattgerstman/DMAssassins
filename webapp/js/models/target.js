// js/models/user.js

var app = app || {Models:{}, Views:{}, Routers:{}, Running:{}, Session:{}};
(function() {
	'use strict';
	
	app.Models.Target = Backbone.Model.extend({
		defaults: {
			'assassin_id' : '',
			'username' : '',
			'user_id' : '',
			'email' : '',
			'properties': {
				'name' : '',
				'facebook': '',
				'twitter': '',
				'photo_thumb' : SPY,
				'photo' : SPY
			}

		},
		parse: function(response) {
                // process response.meta when necessary...
                return response.response;
        },
		urlRoot: config.WEB_ROOT,
		initialize: function() {
			if (!this.get('assassin_id'))
			{
				this.assassin_id = app.Session.get('user_id');
			}
			this.idAttribute = 'assassin_id' 
			this.url = this.urlRoot + app.Session.get('game_id')  + '/users/' + this.assassin_id + '/target/';
		},
		changeUser : function(assassin_id) {
			this.assassin_id = assassin_id;
			this.url = this.urlRoot + app.Session.get('game_id')  + '/users/' + this.assassin_id + '/target/';			
			this.fetch();
		}
	})
})();