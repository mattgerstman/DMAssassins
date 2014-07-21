// js/models/nav.js

var app = app || {};

(function() {
	'use strict';
	
	app.Nav = Backbone.Model.extend({
		defaults: {
			'left' : [
				'Target',
				'My Profile',
				'Leaderboard',
				'Rules'],
			'right' : [			
				{
					'Admin' : [
						'Users',
						'Teams',
						'Plot Twists',
						'Twitter',
						'Game Settings'
					]
				},
				'Logout'
				
			]
		}
		
	})
})();