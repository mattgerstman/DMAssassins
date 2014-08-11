// js/models/user.js

var app = app || {Models:{}, Views:{}, Routers:{}, Running:{}, Session:{}};
(function() {
	'use strict';
	
	app.Models.Profile = Backbone.Model.extend({
		defaults: {
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
			this.idAttribute = 'username'
			this.url = this.urlRoot + this.get('username') + '/';
		},
		changeUser : function(username) {
			this.url = this.urlRoot + username + '/';
			this.fetch();
		}
	})
})();