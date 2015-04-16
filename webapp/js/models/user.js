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
            'user_id': null,
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
        fetch: function(options) {
            if (app.Running.Games.getActiveGameId() === null) {
                return;
            }
            return Backbone.Model.prototype.fetch.call(this, options);
        },
        joinGame: function(game_id, game_password, team_id) {
            var that = this;
            var last_game_id = app.Running.Games.getActiveGameId();
            this.save(null, {
                url: config.WEB_ROOT + 'game/' + game_id + '/user/' + this.get('user_id') + '/',
                type: 'POST',
                data: JSON.stringify({
                    'game_password': game_password,
                    'team_id': team_id
                }),
                success: function() {
                    app.Running.Games.setActiveGame(game_id).set('member', true);
                    that.trigger('join-game');
                },
                error: function(that, response, options) {
                    if (response.status == 401) {
                        that.trigger('join-error-password');
                        app.Running.Games.get(game_id).set('member', false);
                        if (!!last_game_id)
                        {
                            app.Running.Games.setActiveGame(last_game_id, true).set('member', true);
                        }
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
        kill: function(data, successCallback, errorCallback) {
            var that = this;
            var url = this.url() + 'kill/';
            $.ajax({
                url:          url,
                type:         'POST',
                contentType:  'application/json',
                data:         JSON.stringify(data),
                success: function(response) {
                    that.setProperty('alive', 'false');
                    if (typeof successCallback === 'function') {
                        successCallback(response);
                    }
                },
                error: function(response) {
                    if (typeof errorCallback === 'function') {
                        errorCallback(response);
                    }
                }
            });
        },
        revive: function(data, successCallback, errorCallback) {
            var that = this;
            var url = this.url() + 'revive/';
            $.ajax({
                url:          url,
                type:         'POST',
                contentType:  'application/json',
                data:         JSON.stringify(data),
                success: function(response) {
                    that.setProperty('alive', 'true');
                    if (typeof successCallback === 'function') {
                        successCallback(response);
                    }
                },
                error: function(response) {
                    if (typeof errorCallback === 'function') {
                        errorCallback(response);
                    }
                }
            });
        },
        changeRole: function(role_id, options) {
            var url = this.url() + 'role/';
            _.extend(options, {
                type:"PUT",
                url: url,
                data: JSON.stringify({role: role_id}),
            });

            $.ajax(options);
        },
        getRole: function() {

            // get the user role from this user
            var user_role = this.getProperty('user_role');
            if (user_role !== null) {
                return user_role;
            }

            // if we don't have a user role see if this is the same user as the one in the session
            if ((this.get('user_id') === null) || (app.Session.get('user_id') === this.get('user_id'))) {
                return app.Session.get('user_role');
            }

            // if all else fails return null
            return null;
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
        handleRole: function(){
            var user_role = this.getRole();
            app.Session.set('user_role', user_role);
            app.Running.Router.before({}, function(){});
        }
    });
})();
