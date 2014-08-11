// js/models/user.js

var app = app || {Models:{}, Views:{}, Routers:{}, Running:{}, Session:{}};
(function() {
	'use strict';
	
	app.Models.Leaderboard = Backbone.Model.extend({
		defaults: {
			users : [
				{
					name:"Matt",
					kills:75
				},
				{
					name:"Jimmy",
					kills:5
				}
			],
			teams: [
				{
					name: "Slytherin",
					alive: 7,
					kills: 75
				},				
				{
					name: "Ravenclaw",
					alive: 7,
					kills: 14
				},
			],
		},
		url : config.WEB_ROOT+'leaderboard/'

	})
})();