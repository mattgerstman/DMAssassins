// Rules model, loads rules from the db so that admins can define custom rules per game
// js/models/user.js

var app = app || {Models:{}, Views:{}, Routers:{}, Running:{}, Session:{}};
(function() {
	'use strict';
	
	app.Models.Rules = Backbone.Model.extend({
		defaults: {
			rules : {
				'Golden' : ['Don\'t be a dick', 'Don\'t kill Matt'],
				'Safe Zones' : ['No killing on the third floor', 'No killing in Matt\'s house', 'The third floor corridor is off-limits to all those who do not wish to die a most painful death']
			}
		},
		url : config.WEB_ROOT+'rules/'

	})
})();