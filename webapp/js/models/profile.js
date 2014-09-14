// Model for user profile
// js/models/profile.js


(function() {
	'use strict';
	
	app.Models.Profile = Backbone.Model.extend({

		// default profile properties
		defaults: {
			'user_id' : '',
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
				
        // loaded on initialization
		initialize: function() {
			this.idAttribute = 'user_id';
			var game_id = app.Running.Games.getActiveGameId();
			this.url = config.WEB_ROOT + 'game/' + game_id + '/users/' + this.get('user_id') + '/';
		},
		changeGame: function() {
			var game_id = app.Running.Games.getActiveGameId();
			this.url = config.WEB_ROOT + 'game/' + game_id + '/users/' + this.get('user_id') + '/';
		},
		joinGame: function(game_id, game_password) {
    	    this.url = config.WEB_ROOT + 'game/' + game_id + '/users/' + this.get('user_id') + '/';
            var that = this;
    	    this.save(null, {
                headers: {'X-DMAssassins-Game-Password': game_password},
        	    success: function(){        	        
            	    app.Running.Games.addGame(game_id);
            	    that.trigger('join-game');
            	    Backbone.history.navigate('my_profile', { trigger : true });
        	    }
    	    });	
		},
		quit: function(secret) {			
			var that = this;
			this.destroy({
				headers: {'X-DMAssassins-Secret': secret},
				success: function(){
					if (!app.Running.Games.removeActiveGame()) {
						Backbone.history.navigate('#logout', { trigger : true });
						return;
					}
				}				
            })
        },
	})
})();
