// js/models/user.js

var app = app || {Models:{}, Views:{}, Routers:{}, Running:{}};
var logged_in_user = 'Jimmy';
(function() {
	'use strict';
	
	app.Models.User = Backbone.Model.extend({
		defaults: {
			'username' : '',
			'email' : '',
			'properties': {
				'name' : '',
				'facebook': '',
				'twitter': '',
				'photo_thumb' : '',
				'photo' : ''
			}

		},
		parse: function(response) {
                // process response.meta when necessary...
                return response.response;
        },
		urlRoot: WEB_ROOT + 'users/',
		initialize: function() {
			this.url = this.urlRoot + this.get('username') + '/'
		}
	})
})();