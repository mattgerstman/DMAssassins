// js/models/user.js

var app = app || {};

(function() {
	'use strict';
	
	app.User = Backbone.Model.extend({
		defaults: {
			'name' : 'Matt',
			'email' : 'mattgerstman@gmail.com',
			'facebook': 'http://facebook.com',
			'twitter': 'http://twitter.com',
			'photo_thumb' : 'https://graph.facebook.com/jimmyyisflyy/picture?width=300&height=300',
			'photo' : 'https://graph.facebook.com/jimmyyisflyy/picture?width=1000'
		}
	})
})();