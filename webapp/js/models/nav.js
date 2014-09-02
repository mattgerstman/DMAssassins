// js/models/nav.js

var app = app || {Models:{}, Views:{}, Routers:{}, Running:{}, Session:{}};

(function() {
	'use strict';
	
	app.Models.Nav = Backbone.Model.extend({
		
		// Nav setup for top bar, currently there is no serverside representation of this
		defaults: {
			'left' : [
				'Target',
				'My Profile',
				'Leaderboard',
				'Rules'],
			'right' : [			
				{
					'Games' : [
						'Join Another Game'
					],
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
		},
		initialize: function(){
			
		}
		
	})
})();