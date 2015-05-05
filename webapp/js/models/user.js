//
// js/models/user.js
// dmassassins.js
//
// Copyright (c) 2014 Matt Gerstman
// MIT License.
//
// User model, manages single user

(function() {
    'use strict';

    app.Models.User = Backbone.Model.extend({

        // default profile properties
        defaults: {
            'user_id': null,
            'username': null,
            'email': strings.loading,
            'properties': {
                'name': strings.loading,
                'facebook': strings.loading,
                'secret': strings.loading,
                'team': strings.loading,
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
                    if (response.status === 401) {
                        that.trigger('join-error-password');
                        app.Running.Games.get(game_id).set('member', false);
                        if (!!last_game_id)
                        {
                            app.Running.Games.setActiveGame(last_game_id, true).set('member', true);
                        }
                    }
                }
            });
            return this;
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
            return this;
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
            return this;
        },
        changeRole: function(role_id, options) {
            this.set('role', role_id);
            this.setProperty('user_role', role_id);
            var url = this.url() + 'role/';
            options.url = url;
            return this.save(null, options);
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
        changeTeam: function(team_id, team_name, success, error) {
            var that = this;
            return this.save(null, {
                url: this.url() + 'team/' + team_id + '/',
                success: function(user, response) {
                    user.setProperty('team', team_name);
                    user.set('team', team_id);
                    console.log(user);
                    if (typeof success === 'function') {
                        success(user, response);
                    }
                },
                error: function(response, user) {
                    if (typeof error === 'function') {
                        error(user, response);
                    }
                }
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
            return this;
        },
        handleRole: function(){
            var user_role = this.getRole();
            app.Session.set('user_role', user_role);
            app.Running.Router.before({}, function(){});
            return this;
        }
    });
})();
