// js/models/user.js

var app = app || {Models:{}, Views:{}, Routers:{}, Running:{}, Session:{}};
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
			var trailing = this.get('type') == 'target' ? 'target/' : '';
			this.url = this.urlRoot + this.get('username') + '/' + trailing;
		},
		changeUser : function(username) {
			var trailing = this.get('type') == 'target' ? 'target/' : '';
			this.url = this.urlRoot + username + '/' + trailing;
			this.fetch();
		}
	})
})();