// js/models/user.js

var app = app || {};

(function() {
	'use strict';
	
	app.User = Backbone.Model.extend({
		defaults: {
			'email' : 'Matt',
			'properties': {
				'name' : 'Jimmy',
				'facebook': 'http://facebook.com',
				'twitter': 'http://twitter.com',
				'photo_thumb' : 'https://graph.facebook.com/jimmyyisflyy/picture?width=300&height=300',
				'photo' : 'https://graph.facebook.com/jimmyyisflyy/picture?width=1000'
			}
		},
		parse: function(response) {
                // process response.meta when necessary...
                return response.response;
        },
		urlRoot: WEB_ROOT + 'users/',
		initialize: function() {
			this.url = this.urlRoot + this.get('email')
			this.fetch()
		}
	})
})();