//
// js/models/user.js
// dmassassins.js
//
// Copyright (c) 2014 Matt Gerstman
// MIT License.
//
// User model, manages single user

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

    app.Models.User = Backbone.Model.extend({

        // default profile properties
        defaults: {
            'user_id': '',
            'username': '',
            'email': 'Loading...',
            'properties': {
                'name': 'Loading..',
                'facebook': 'Loading..',
                'secret': 'Loading..',
                'team': 'Loading..',
                'photo_thumb': SPY,
                'photo': SPY
            }

        },
        idAttribute : 'user_id',
        url: function() {
            var game_id = app.Running.Games.getActiveGameId();
            return config.WEB_ROOT + 'game/' + game_id + '/user/' + this.get('user_id') + '/';           
        },       
        joinGame: function(game_id, game_password, team_id) {
            var that = this;            
            var last_game_id = app.Running.Games.getActiveGameId();
            this.save(null, {
                url: config.WEB_ROOT + 'game/' + game_id + '/user/' + this.get('user_id') + '/',
                headers: {
                    'X-DMAssassins-Game-Password': game_password,
                    'X-DMAssassins-Team-Id': team_id
                },
                success: function() {
                    app.Running.Games.setActiveGame(game_id).set('member', true);                    
                    that.trigger('join-game');
                },
                error: function(that, response, options) {
                    if (response.status == 401) {
                        that.trigger('join-error-password');
                        app.Running.Games.get(game_id).set('member', false);
                        app.Running.Games.setActiveGame(last_game_id, true).set('member', true);
                    }
                }
            });
        },
        setProperty: function(key, value, silent) {
            var properties = this.get('properties');
            if (!properties)
                properties = {};
            properties[key] = value;
            this.set('properties', properties);
            if ((silent === undefined) || (silent === false))
            {
                this.trigger('change');
            }            
            return this.get('properties');
        },
        getProperty: function(key){
            var properties = this.get('properties');
            if (!properties)
                return null;
            if (properties[key] === undefined)
                return null;
            return properties[key];
        },
        kill: function(){
        	var that = this;
        	var url = this.url() + 'kill/';
	      	$.post(url, function(response){
	      		that.setProperty('alive', 'false');
	      	});
        },
        revive: function(){
        	var that = this;
        	var url = this.url() + 'revive/';
	      	$.post(url, function(response){
	      		that.setProperty('alive', 'true');
	      	});
        },
        saveEmailSettings: function(email, allow_email) {
            var data = {
                email: email,
                allow_email: allow_email
            };
            var that = this;
            var url = this.url() + 'email/';
            $.post(url, data, function(response){
                that.set('email', email);
            });
        },
        quit: function(secret) {
            var that = this;
            this.destroy({
                headers: {
                    'X-DMAssassins-Secret': secret
                },
                success: function() {
                    if (!app.Running.Games.removeActiveGame()) {
                        Backbone.history.navigate('#logout', {
                            trigger: true
                        });
                        return;
                    }
                },
                error: function(model, response){
                    alert(response.responseText);
                }
            });
        },
        checkAccess: function(){
            app.Running.Router.before({}, function(){});
        }
    });
})();
