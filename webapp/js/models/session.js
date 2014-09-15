//
// js/models/session.js
// dmassassins.js
//
// Copyright (c) 2014 Matt Gerstman
// MIT License.
//
// Manages all local storage information and helps keep various models in sync

var app = app || {
    Collections: {},
    Models: {},
    Views: {},
    Routers: {},
    Running: {},
    Session: {}
};

(function() {

    app.Models.Session = Backbone.Model.extend({

        url: config.WEB_ROOT + 'session/',

        initialize: function() {

            // Check for localstorage support
            if (Storage && sessionStorage) {
                this.supportStorage = true;
            }
        },

        // returns data stored in the session
        get: function(key) {
            if (this.supportStorage) {
                var data = sessionStorage.getItem(key);
                if (data && data[0] === '{') {
                    return JSON.parse(data);
                } else {
                    return data;
                }
            } else {
                return Backbone.Model.prototype.get.call(this, key);
            }
        },

        // sets a session variable
        set: function(key, value) {
            if (this.supportStorage) {
                sessionStorage.setItem(key, value);
            } else {
                Backbone.Model.prototype.set.call(this, key, value);
            }
            return this;
        },

        // unsets a session 
        unset: function(key) {
            if (this.supportStorage) {
                sessionStorage.removeItem(key);
            } else {
                Backbone.Model.prototype.unset.call(this, key);
            }
            return this;
        },

        // clears all data from the session
        clear: function() {
            if (this.supportStorage) {
                sessionStorage.clear();
            } else {
                Backbone.Model.prototype.clear(this);
            }
        },
        // calls the facebook login function and handles it appropriately
        // if they are logged into facebook and connected to the app a session is created automatically
        // otherwise a popup will appear and handle the session situation
        login: function() {

            var parent = this;

            FB.getLoginStatus(function(response) {
                if (response.status === 'connected') {
                    // Logged into your app and Facebook.
                    //console.log(response);
                    parent.createSession(response);


                } else if (response.status === 'not_authorized') {

                    // The person is logged into Facebook, but not your app.
                    FB.login(function(response) {
                        parent.createSession(response);

                        // scope are the facebook permissions we're requesting 
                    }, {
                        scope: 'public_profile,email,user_friends,user_photos'
                    })

                } else {

                    FB.login(function(response) {
                        parent.createSession(response);

                        // scope are the facebook permissions we're requesting
                    }, {
                        scope: 'public_profile,email,user_friends,user_photos'
                    })

                    // The person is not logged into Facebook, so we're not sure if
                    // they are logged into this app or not.
                }
            })

        },

        // takes a facebook response and creates a session from it
        createSession: function(response) {

            var game_id = this.get('game_id');

            var data = {
                'facebook_id': response.authResponse.userID,
                'facebook_token': response.authResponse.accessToken,
                'game_id': game_id
            }

            var that = this;

            // performs the ajax request to the server to get session data
            var login = $.ajax({
                url: this.url,
                data: data,
                type: 'POST'
            });

            // after the ajax request run this function
            login.done(function(response) {

                // store all reponse data in the new session immediately 
                var parsedGames    = app.Running.Games.parse(response.games);    
                var games 		   = parsedGames           || {};
                var game 		   = response.game         || { game_id: null };
                var user 		   = response.user         || { user_id: null };
                var target 		   = response.target       || { assassin_id: user.user_id };                
                var leaderboard    = response.leaderboard  || {};        
                var rules          = null;
                if (game.game_properties)
                {
                    rules = {rules: (game.game_properties.rules || null)};
                }

                target.assassin_id = response.user.user_id;

                // reload the data for all models
                app.Running.ProfileModel.set(user);
                app.Running.TargetModel.set(target);               
                app.Running.LeaderboardModel.set(leaderboard);
                app.Running.RulesModel.set(rules);
                app.Running.Games.reset(games);
                
                // store the basic auth token in the session in case we need to reload it on app launch
                that.storeSession(response)
                app.Running.Games.setActiveGame(game.game_id, true);

                var targetURLs = app.Running.Router.requiresTarget;
                var path = Backbone.history.fragment;

                if ((path != 'login') && !_.contains(targetURLs, path)) {
                    Backbone.history.loadUrl();
                    return;
                }

                if (path == 'login' && app.Running.TargetModel.get('user_id')) {
                    Backbone.history.navigate('#', {
                        trigger: true
                    });       
                    return;
                }

                Backbone.history.navigate('#my_profile', {
                    trigger: true
                });                
                


            });

            // if theres a login error direct them to the login screen
            login.fail(function() {
                that.clear();
                FB.getLoginStatus(function(response) {
                    statusChangeCallback(response);
                });
                Backbone.history.navigate('login', {
                    trigger: true
                });
            });


        },
        storeSession: function(data) {
            // store a boolean to determine if we're authenticated
            this.set('authenticated', true); 
            this.set('has_game', data.game !== null);
            this.storeBasicAuth(data);
        },
        // stores all the basic auth variables in the session
        storeBasicAuth: function(data) {
            var user_id = data.user.user_id;
            this.set('user_id', user_id)

            var token = data.token
            var plainKey = user_id + ":" + token
            var base64Key = window.btoa(plainKey);
            this.set('authKey', base64Key);
            this.setAuthHeader();
        },

        // sets the Basic Auth header for all ajax requests
        setAuthHeader: function() {
            var base64Key = this.get('authKey');
            $.ajaxSetup({
                headers: {
                    'Authorization': "Basic " + base64Key
                }
            });

        }
    });
})()