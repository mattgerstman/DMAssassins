//
// js/models/session.js
// dmassassins.js
//
// Copyright (c) 2014 Matt Gerstman
// MIT License.
//
// Manages all local storage information and helps keep various models in sync

(function() {

    app.Models.Session = Backbone.Model.extend({

        url: config.WEB_ROOT + 'session/',

        initialize: function() {
            try {
                // Check for localstorage support
                if (Storage && localStorage) {
                    this.supportStorage = true;
                }
                if (!this.supportStorage) {
                    return;
                }
                localStorage.setItem('test', 0);
            }
            catch (err) {
                this.supportStorage = false;
            }
        },

        // returns data stored in the session
        get: function(key) {
            if (this.supportStorage) {
                var data = localStorage.getItem(key);
                try {
                    return JSON.parse(data);
                } catch (err) {
                    console.log("ERROR ACCESSING SESSION DATA");
                    console.log(err);
                    console.log(data);
                    return data;
                }
            } else {
                return Backbone.Model.prototype.get.call(this, key);
            }
        },

        // sets a session variable
        set: function(key, value) {
            if (this.supportStorage) {
                localStorage.setItem(key, JSON.stringify(value));
            } else {
                Backbone.Model.prototype.set.call(this, key, value);
            }
            return this;
        },

        // unsets a session
        unset: function(key) {
            if (this.supportStorage) {
                localStorage.removeItem(key);
            } else {
                Backbone.Model.prototype.unset.call(this, key);
            }
            return this;
        },

        // clears all data from the session
        clear: function() {
            if (this.supportStorage) {
                localStorage.clear();
            } else {
                Backbone.Model.prototype.clear(this);
            }
            return this;
        },
        // logs the user out
        logout: function() {
            var options = {silent:true};

            this.clear();
            app.Running.Games.clearActiveGame(options);
            app.Running.Users.reset(options);
            app.Running.Teams.reset(options);

            app.Running.User.clear(options);
            app.Running.Permissions.clear(options);
            app.Running.TargetModel.clear(options);
            app.Running.LeaderboardModel.clear(options);
            app.Running.RulesModel.clear(options);
            app.Running.TargetFriendsModel.clear(options);

        },
        // calls the facebook login function and handles it appropriately
        // if they are logged into facebook and connected to the app a session is created automatically
        // otherwise a popup will appear and handle the session situation
        login: function() {

            var parent = this;

            try
            {
                if (!app.Running.FB) {
                    throw new Error ('Error loading Facebook SDK');
                }
            }
            catch (e)
            {
                Raven.captureException(e);
                alert('There was an issue connecting to Facebook, please refresh and try again in a minute.');
            }

            app.Running.FB.getLoginStatus(function(response) {
                if (response.status === 'connected') {
                    // Logged into your app and Facebook.
                    //console.log(response);

                    try
                    {
                        if (!response.authResponse)
                            throw new Error('Error processing facebook login');

                        parent.createSession(response);
                    }
                    catch (e)
                    {
                        Raven.captureException(e, {extra: response});
                        alert('Your session has expired. Please log in again.');
                    }

                } else if (response.status === 'not_authorized') {

                    // The person is logged into Facebook, but not your app.
                    app.Running.FB.login(function(response) {

                        if (response.authResponse)
                        {
                            parent.createSession(response);
                        }
                        else
                        {
                            alert('You must authorize Facebook to play DMAssassins!');
                            location.reload();
                        }
                    }, {
                        // scope are the facebook permissions we're requesting
                        scope: config.ALL_PERMISSIONS
                    });

                } else {
                    // The person is not logged into Facebook, so we're not sure if
                    // they are logged into this app or not.

                    // hack for chrome for iOS which doesn't like the native FB redirect
                    if( navigator.userAgent.match('CriOS') )
                    {
                        window.open('https://www.facebook.com/dialog/oauth?client_id='+config.APP_ID+'&redirect_uri='+ config.CLIENT_ROOT+'%23login&scope='+config.ALL_PERMISSIONS, '', null);
                        return;
                    }
                    app.Running.FB.login(function(response) {

                        if (response.authResponse)
                        {
                            parent.createSession(response);
                        }
                        else
                        {
                            alert('You must authorize Facebook to play DMAssassins!');
                            location.reload();
                        }
                    }, {
                        // scope are the facebook permissions we're requesting
                        scope: config.ALL_PERMISSIONS
                    });
                }
                return this;
            });

        },
        recoverSession: function(response) {
            if (response.status === 'connected') {
                // Logged into your app and Facebook.
                try
                {
                    if (!response.authResponse)
                        throw new Error('Error processing facebook login');

                    return this.createSession(response);
                }
                catch (e)
                {
                    Raven.captureException(e, {extra: response});
                }
            }
            this.clear();
            Backbone.history.navigate('', {
                trigger: true
            });
        },
        handleResponse: function(response) {

                // Get permissions from facebook asynchronously
                app.Running.Permissions.fetch();

                // store all reponse data in the new session immediately
                var parsedGames    = app.Running.Games.parse(response.games);
                var game 		       = response.game         || { game_id: null };
                var user 		       = response.user         || { user_id: null };


                // Set user id in case an error occurs
                Raven.setUser({
                    user_id: user.user_id,
                    game_id: game.game_id
                });

                if (!response.token) {
                    Raven.captureException(new Error("Server didn't return token"), {extra: user});
                    alert('An unexpected error occurred. Please try again');
                    app.Session.clear();
                    Backbone.history.navigate('', {
                        trigger: true
                    });
                    return;
                }

                // reload the data for all models
                app.Running.User.set(user);

                var last_role = app.Session.get('user_role');
                app.Running.User.setProperty('user_role', last_role);

                // store the basic auth token in the session in case we need to reload it on app launch
                app.Session.storeSession(response);
                if (game.game_id) {
                    app.Running.Games.add(game);
                    app.Running.Games.setActiveGame(game.game_id);
                    app.Running.Games.fetch({reset:true});
                }

                var targetURLs = app.Running.Router.requiresTarget;
                var path = Backbone.history.fragment;

                if ((path !== 'login') && !_.contains(targetURLs, path)) {
                    Backbone.history.loadUrl();
                    return;
                }

                if (path === 'login' && app.Running.TargetModel.get('user_id')) {
                    Backbone.history.navigate('#target', {
                        trigger: true
                    });
                    return;
                }

                Backbone.history.navigate('#my-profile', {
                    trigger: true
                });

        },
        // takes a facebook response and creates a session from it
        createSession: function(response) {

            var game_id = this.get('game_id');

            var data = {
                'facebook_id': response.authResponse.userID,
                'facebook_token': response.authResponse.accessToken,
                'game_id': game_id
            };

            var that = this;

            // performs the ajax request to the server to get session data
            var login = $.ajax({
                url: this.url,
                data: JSON.stringify(data),
                contentType: 'application/json',
                tryCount:0,
                retryLimit:3,
                type: 'POST',
                success: that.handleResponse,
                error: function(serverResponse, textStatus, errorThrown) {
                    Raven.captureException(new Error('Server failed to login'), {extra: {facebook_response:response, try_count: this.tryCount, server_response: serverResponse, text_status: textStatus, error_thrown :errorThrown}});
                    that.clear();

                    alert('An error occurred. Please try again later.');
                    Backbone.history.navigate('', {
                        trigger: true
                    });
                }
            });
        },
        // store a boolean to determine if we're authenticated
        storeSession: function(data) {
            this.set('authenticated', true);
            this.set('has_game', data.game !== null);
            this.storeBasicAuth(data);
            return this;
        },
        // stores all the basic auth variables in the session
        storeBasicAuth: function(data) {
            var user_id = data.user.user_id;
            this.set('user_id', user_id);

            var token = data.token;
            var plainKey = user_id + ":" + token;
            var base64Key = window.btoa(plainKey);
            this.set('authKey', base64Key);
            this.setAuthHeader();
            return this;
        },
        // sets the Basic Auth header for all ajax requests
        setAuthHeader: function() {
            var base64Key = this.get('authKey');
            $.ajaxSetup({
                headers: {
                    'Authorization': "Basic " + base64Key
                }
            });
            return this;
        }
    });
})();
