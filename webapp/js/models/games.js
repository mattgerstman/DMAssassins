// js/models/games.js

var app = app || {Models:{}, Views:{}, Routers:{}, Running:{}, Session:{}};
(function() {
	'use strict';
	
	app.Models.Games = Backbone.Model.extend({
		defaults: {
			games: [
				{
					game_id : '',
					game_name : 'Triwizard Tournament',
					game_started : false
				}
			]
		},
		parse: function(response) {
            // process response.meta when necessary...
            var wrapper = {
                games: response.response
            };
            
            return wrapper;
        },
		url : config.WEB_ROOT+'game/'

	})
})();