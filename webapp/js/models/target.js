// model for target pages
// js/models/user.js

var app = app || {
    Collections: {},
    Models: {},
    Views: {},
    Routers: {},
    Running: {},
    Session: {}
};

(function() {
	'use strict';
	
	app.Models.Target = Backbone.Model.extend({
		defaults: {
		    'game_id' : null,
			'assassin_id' : '',
			'username' : '',
			'user_id' : '',
			'email' : '',
			'properties': {
				'name' : '',
				'facebook': '',
				'twitter': '',
				'photo_thumb' : SPY,
				'photo' : SPY
			}

		},
		// called automatically by fetch() to remove the wrapper
		parse: function(response) {
			return response;
        },       
        // consstructor
		initialize: function() {
			if (!this.get('assassin_id'))
			{
				this.assassin_id = app.Session.get('user_id');
			}
			this.idAttribute = 'assassin_id'
			var game_id = app.Running.Games.getActiveGameId();
			this.url = config.WEB_ROOT + "game/" + game_id  + '/users/' + this.get('assassin_id') + '/target/';
		},
		
		// change the user who's target we're getting
		changeGame : function(game_id) {
		    if (this.get('game_id') != game_id)
		    {
		        this.set('game_id', game_id);
    		    this.url = config.WEB_ROOT + "game/" + game_id  + '/users/' + this.get('assassin_id') + '/target/';    
    		    return true;
		    }
		    return false;
			
		}
	})
})();
