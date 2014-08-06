// js/models/user.js

var app = app || {Models:{}, Views:{}, Routers:{}, Running:{}, Session:{}};
(function() {
	'use strict';
	
	app.Models.Target = Backbone.Model.extend({
		defaults: {
			'assassin' : '',
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
		parse: function(response) {
                // process response.meta when necessary...
                return response.response;
        },
		urlRoot: config.WEB_ROOT + 'users/',
		initialize: function() {
			this.idAttribute = 'assassin' 
			this.url = this.urlRoot + this.get('assassin') + '/target/';
		},
		changeUser : function(username) {
			var trailing = this.get('type') == 'target' ? 'target/' : '';
			this.url = this.urlRoot + username + '/' + trailing;
			this.fetch();
		}
	})
})();