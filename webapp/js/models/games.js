// js/models/games.js

var app = app || {Models:{}, Views:{}, Routers:{}, Running:{}, Session:{}};
(function() {
	'use strict';
	
	app.Models.Games = Backbone.Model.extend({
		defaults: {
			user_id: null,
			games: [
				{
					game_id : '',
					game_name : 'Triwizard Tournament',
					game_started : false,
					game_has_password : false
				}
			]
		},
		initialize: function(){
			this.url = config.WEB_ROOT+'game/'+this.user_id	
		},
		parse: function(response) {
            // process response.meta when necessary...
            var wrapper = {
                games: response.response
            };
            
            return wrapper;
        },
        create: function(name, password){
	        var that = this;
	        var data = {
        		user_id: app.Session.get('user_id'),
	        	game_name: name,
	        	game_password: password
	        	
        	}
        	$.post(this.url, data, function(response){
		        app.Session.setGame(response.response)
		        that.trigger('finish_set_game');
	        })
        },
        join: function(game_id, user_id, password){
        	var that = this;
        	var data = {
	        	game_password: password    	
        	}
	        $.post(config.WEB_ROOT+game_id+'/users/'+user_id+'/', data, function(response){
		        app.Session.setGame(response.response)
		        that.trigger('finish_set_game');
	        })
        },
        switchGame: function(game_id){
        	var that = this;
	      $.get(config.WEB_ROOT+'game/'+game_id+'/', function(response){
		      app.Session.setGame(response.response);
		      that.trigger('game-change');
	      })  
        },
        loadUser: function(user_id){
	        this.user_id = user_id;
	        this.url = config.WEB_ROOT+'users/'+this.user_id+'/game/';
        }        
		

	})
})();