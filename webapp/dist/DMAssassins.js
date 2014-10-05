/*! DMAssassins - v0.9.0 - 2014-10-05
* http://dmassassins.com
* Copyright (c) 2014 Matt Gerstman; Licensed  */
// Used to handle user permissions
// All APIs should be equipped with the right set of permissions so this simply prevents users from stumbling onto a page they can't use

var AuthUtils = {
    getRolesMap: function() {
        var rolesMap = {
            dm_user:        {value: 0, pretty_name: "User"},
            dm_captain:     {value: 1, pretty_name: "Captain"},
            dm_admin:       {value: 2, pretty_name: "Admin"},
            dm_super_admin: {value: 3, pretty_name: "Super Admin"},
        }
        return rolesMap;    
    },
    getRolesMapFor: function(role, teams_enabled) {
        if (!role) {
            return {};
        }
        var filteredMap = {};
        var rolesMap = this.getRolesMap();
        for (var key in rolesMap)  {
            if (teams_enabled === false) {
                if (rolesMap[key]['value'] == 1)
                    continue;
            }
            if (rolesMap[key]['value'] <= rolesMap[role]['value'])
                filteredMap[key] = rolesMap[key];
        }
        return filteredMap;
    },
    getRolePrettyName: function(role) {
        var rolesMap = this.getRolesMap();
        return rolesMap[role]['pretty_name'];
    },
    isRoleAllowed: function(userRole, minRoleForAccess) {
        var rolesMap = this.getRolesMap();
        if ((userRole === undefined) || (rolesMap[userRole] === undefined))
            return undefined
        
        return rolesMap[userRole]['value'] >= rolesMap[minRoleForAccess]['value'];        
    },
    requiresUser: function(userRole) {
        return this.isRoleAllowed(userRole, 'dm_user');
    },
    requiresCaptain: function(userRole) {
        return this.isRoleAllowed(userRole, 'dm_captain');
    },
    requiresAdmin: function(userRole) {
        return this.isRoleAllowed(userRole, 'dm_admin');
    },
    requiresSuperAdmin: function(userRole) {
        return this.isRoleAllowed(userRole, 'dm_super_admin');
    }
}


var SPY = '/assets/img/spy.jpg';

var app = app || {
    Collections: {},
    Models: {},
    Views: {},
    Routers: {},
    Running: {},
    Session: {}
};




  // This is called with the results from from FB.getLoginStatus().
  function statusChangeCallback(response) {
    // The response object is returned with a status field that lets the
    // app know the current login status of the person.
    // Full docs on the response object can be found in the documentation
    // for FB.getLoginStatus().
    
    if (!app.Session.get('authenticated'))
    {
    	//console.log('no facebook response');
	    return;
    }
    	
    
    if (response.status === 'connected') {
      // Logged into your app and Facebook.

		  //console.log('connected');
		  app.Session.createSession(response, function(){
//	  		  app.Running.Router.reload();
				
		  });


    } else {
		//console.log('else');

		app.Running.Router.navigate('login')
    }
    
  }




  window.fbAsyncInit = function() {
  FB.init({
    appId      : config.APP_ID,
    cookie     : true,  // enable cookies to allow the server to access 
                        // the session
    xfbml      : true,  // parse social plugins on this page
    version    : 'v2.0', // use version 2.0
    status	   : true
  });
  
  app.Running.FB = FB;

  // Now that we've initialized the JavaScript SDK, we call 
  // FB.getLoginStatus().  This function gets the state of the
  // person visiting this page and can return one of three states to
  // the callback you provide.  They can be:
  //
  // 1. Logged into your app ('connected')
  // 2. Logged into Facebook, but not your app ('not_authorized')
  // 3. Not logged into Facebook and can't tell if they are logged into
  //    your app or not.
  //
  // These three cases are handled in the callback function.


	  FB.getLoginStatus(function(response) {
	    statusChangeCallback(response);
	  });


  };

  // Load the SDK asynchronously
  (function(d, s, id) {
    var js, fjs = d.getElementsByTagName(s)[0];
    if (d.getElementById(id)) return;
    js = d.createElement(s); js.id = id;
	js.src = "//connect.facebook.net/en_US/sdk.js";
//    js.src = "//connect.facebook.net/en_US/sdk/debug.js";
    fjs.parentNode.insertBefore(js, fjs);
  }(document, 'script', 'facebook-jssdk'));

//
// js/models/games.js
// dmassassins.js
//
// Copyright (c) 2014 Matt Gerstman
// MIT License.
//
// Single game model

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
    app.Models.Game = Backbone.Model.extend({

        // default properties with a fake game
        defaults: {
            game_id: '',
            game_name: 'Loading...',
            game_started: false,
            game_has_password: false,
            member: true
        },

        idAttribute: 'game_id',
        urlRoot: config.WEB_ROOT + 'user/',
        areTeamsEnabled: function() {
            return this.getProperty('teams_enabled') == 'true';
        },
        getProperty: function(key) {
            var properties = this.get('game_properties');
            if (!properties)
                return null;
            if (properties[key] === undefined)
                return null;
            return properties[key]  
        },
        fetchProperties: function() {
            var url = this.gameUrl();
            return this.fetch({url: url});
        },
        gameUrl: function() {
            var url = config.WEB_ROOT + 'game/' + this.get('game_id') + '/';
            return url;
        },
        url: function() {
            var url = this.urlRoot;
            url += app.Session.get('user_id') + '/game/';            
            var game_id = this.get('game_id');
            if (!game_id)
            {
                return url;
            }            
            return url + game_id + '/';
        }
    })
})();

//
// js/models/leaderboard.js
// dmassassins.js
//
// Copyright (c) 2014 Matt Gerstman
// MIT License.
//
// Leaderboard model, displays high scores for a game

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

    app.Models.Leaderboard = Backbone.Model.extend({
        defaults: {
            teams_enabled: true,
            user_col_width: 20,
            team_col_width: 20,
            users: [{
                name: "Loading...",
                kills: "Loading...",
                team_name: "Loading..."
            }, {
                name: "Loading...",
                kills: "Loading...",
                team_name: "Loading..."
            }],
            teams: [{
                "Loading...": "Loading..."
            }]
        },
        parse: function(data) {
            data.user_col_width = 100 / (3 + this.get('teams_enabled'));
            data.team_col_width = 20;
            return data;
        },
        url: function() {
            var game_id = app.Running.Games.getActiveGameId();
            return config.WEB_ROOT + 'game/' + game_id + '/leaderboard/'
        }
    })
})();
//
// js/models/user.js
// dmassassins.js
//
// Copyright (c) 2014 Matt Gerstman
// MIT License.
//
// Rules model, loads rules from the db so that admins can define custom rules per game


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

    app.Models.Rules = Backbone.Model.extend({
        defaults: {
            rules: 'Loading...'
        },
        url: function() {
            var game_id = app.Running.Games.getActiveGameId();
            return config.WEB_ROOT + 'game/' + game_id + '/rules/';
        }
    })
})();
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
                app.Running.User.set(user);
                app.Running.TargetModel.set(target);               
                app.Running.LeaderboardModel.set(leaderboard);
                app.Running.RulesModel.set(rules);
                app.Running.Games.reset(games);
                
                // store the basic auth token in the session in case we need to reload it on app launch
                that.storeSession(response)
                app.Running.Games.setActiveGame(game.game_id, true);
                app.Running.Games.getActiveGame().set(game);

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
//
// js/models/target.js
// dmassassins.js
//
// Copyright (c) 2014 Matt Gerstman
// MIT License.
//
// model for target pages

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
            'game_id': null,
            'assassin_id': '',
            'username': '',
            'user_id': '',
            'properties': {
                'name': 'Loading...',
                'facebook': 'Loading...',
                'team':'Loading...',
                'photo_thumb': SPY,
                'photo': SPY
            }
        },
        url: function() {
            var game_id = app.Running.Games.getActiveGameId();
            return config.WEB_ROOT + "game/" + game_id + '/user/' + this.get('assassin_id') + '/target/';
        },
        // consstructor
        initialize: function() {
            if (!this.get('assassin_id')) {
                this.assassin_id = app.Session.get('user_id');
            }
            this.idAttribute = 'assassin_id'
        }
    })
})();

//
// js/models/team.js
// dmassassins.js
//
// Copyright (c) 2014 Matt Gerstman
// MIT License.
//
// Team model, manages single team

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

    app.Models.Team = Backbone.Model.extend({

        // default profile properties
        defaults: {
            'user_id': '',
            'team_id': '',
            'team_name': ''
        },
        idAttribute : 'team_id',
        url: function(){            
            var game_id = app.Running.Games.getActiveGameId();
            var user_id = this.get('user_id');
            var team_id = this.get('team_id');
           
            if (!!user_id) {
                return config.WEB_ROOT + 'game/' + game_id + '/user/' + user_id + '/team/' + team_id + '/';
            }
            return config.WEB_ROOT + 'game/' + game_id + '/team/' + team_id + '/';
        }

    })
})();

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
            app.Running.Games.setActiveGame(game_id).set('member', true);
            this.save(null, {
                headers: {
                    'X-DMAssassins-Game-Password': game_password,
                    'X-DMAssassins-Team-Id': team_id
                },
                success: function() {
                    that.trigger('join-game');
                    
                    
                    Backbone.history.navigate('my_profile', {
                        trigger: true
                    });
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
        ban: function(){
        	var that = this;
	      	this.destroy({
		      	url: that.url() + 'ban/'
	      	})  
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
            }
            var that = this;
            var url = this.url() + 'email/';
            $.post(url, data, function(response){
                that.set('email', email);
            })
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
                }
            })
        },
        checkAccess: function(){
            app.Running.Router.before({}, function(){});
        }
    })
})();

// Games Collection, Handles all of the games a user has access to
// js/collections/games.js
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
    app.Collections.Games = Backbone.Collection.extend({

        model: app.Models.Game,

        parse: function(response) {
            if (!response) {
                return null;
            }
            var list = [];
            _.each(response.available, function(item){
                item.member = false;
                list.push(item);
            })
            _.each(response.member, function(item){
                item.member = true;
                list.push(item);
            })
            return list;
        },
        comparator: 'game_name',
        active_game: null,
        // handle on initiliazation
        url:function(){
            var user_id = app.Session.get('user_id');
            return config.WEB_ROOT + 'user/' + user_id + '/game/';
            
        },
        addGame: function(game_id){
            var that = this;
            var game = new app.Models.Game({game_id: game_id});
            game.fetch({success: function(){
                that.add(game);
                that.setActiveGame(game.get('game_id'));
            }})
        },
        joinGame: function(game_id, password, team_id) {
            app.Running.User.joinGame(game_id, password, team_id);  
            this.trigger('game-change');          
        },
        setArbitraryActiveGame: function(silent) {
            var newGame = this.findWhere({game_started: true})
            if (!newGame)
            {
                newGame = this.findWhere({game_started: false})
            }
            this.setActiveGame(newGame, silent);
            return newGame;
        },
        removeActiveGame: function(){
            this.remove(this.active_game);
            return this.setArbitraryActiveGame();
        },
        setActiveGame: function(game_id, silent) {
            var game = this.get(game_id);
            if (!game) {
                return null;
            }
            game.fetchProperties();
            this.active_game = game;
            app.Session.set('game_id', game_id);
            app.Session.set('has_game', true);

            if (silent === undefined || !silent)
            {   
                this.trigger('game-change');    
            }
            
            return this.active_game;
        },
        getActiveGame: function() {    
            if (!this.active_game)
            {
                this.setArbitraryActiveGame(true);
            }    
            return this.active_game;

        },
        getActiveGameId: function() { 
            var game = this.getActiveGame();
            if (!game)
            {
                return null;
            }
            return game.get('game_id');
        }

        
        
    })
})();

// Users Collection. Handles all of the users for a game for an admin
// js/collections/users.js
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
    app.Collections.Teams = Backbone.Collection.extend({

        model: app.Models.Team,
        parse: function(response){
            return _.values(response);  
        },
        url: function(){
            var game_id = app.Running.Games.getActiveGameId();
            return config.WEB_ROOT + 'game/' + game_id + '/team/';
        }
        
    })
})();

// Users Collection. Handles all of the users for a game for an admin
// js/collections/users.js
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
    app.Collections.Users = Backbone.Collection.extend({

        model: app.Models.User,
        parse: function(response){
            return _.values(response);  
        },
        url: function(){
            var game_id = app.Running.Games.getActiveGameId();
            return config.WEB_ROOT + 'game/' + game_id + '/users/';
        }
        
    })
})();

//
// js/views/admin-edit-rules-view.js
// dmassassins.js
//
// Copyright (c) 2014 Matt Gerstman
// MIT License.
//
// displays rules for a game


var app = app || {
    Collections: {},
    Models: {},
    Views: {},
    Routers: {},
    Running: {},
    Session: {}
};

(function($) {
    'use strict';
    app.Views.AdminEditRulesView = Backbone.View.extend({


        template: _.template($('#admin-edit-rules-template').html()),
        tagName: 'div',

        initialize: function(params) {
            this.model = app.Running.RulesModel;

            this.listenTo(this.model, 'change', this.render)
            this.listenTo(this.model, 'fetch', this.render)
            this.listenTo(this.model, 'set', this.render)
        },
        loadEditor: function(){
            var that = this;
            this.$el.find("#rules-editor").markdown({
                savable:true,
                saveButtonClass: 'btn btn-md btn-primary',
                footer: '<div class="saved hide">Saving...</div>',
                onSave: function(event) {
                        var rules = event.getContent();
                        that.model.set('rules', rules);
                        $('.saved').removeClass('hide')
                        that.model.save(null, {success: function(){
                            $('.saved').text('Saved.').fadeOut(2000, function(){
                                $(this).text('Saving...');    
                            });
                            
                        }});
                    },
                })
        },
        render: function() {
            var data = this.model.attributes;
            this.$el.html(this.template(data));         
            this.loadEditor();
            return this;
        }

    })

})(jQuery);
//
// js/views/admin-user-view.js
// dmassassins.js
//
// Copyright (c) 2014 Matt Gerstman
// MIT License.
//
// displays user profile


var app = app || {
    Collections: {},
    Models: {},
    Views: {},
    Routers: {},
    Running: {},
    Session: {}
};

(function($) {
    'use strict';
    app.Views.AdminGameSettingsView = Backbone.View.extend({

        template: _.template($('#admin-game-settings-template').html()),
        tagName:'div',
        events: {
            'click .save-game': 'saveGame',
            'click .start-game': 'startGameModal',
            'click .start-game-submit': 'startGame',
            'click .end-game': 'endGameModal',
            'click .end-game-submit': 'endGame'

        },
        initialize: function(){
            this.model = app.Running.Games.getActiveGame();
            this.listenTo(this.model, 'fetch', this.render)
            this.listenTo(this.model, 'change', this.render)
            this.listenTo(this.model, 'save', this.render)
        },
        saveGame: function(event){
        
            // Get values from form
            var game_name = $('#game_name').val();
            var game_password = $('#game_password').val();
            var game_teams_enabled = $('#teams_enabled').is(':checked') ? 'true' : 'false';
            
            // Set values in model
            this.model.set({
                game_name: game_name,
                game_password: game_password,
                game_teams_enabled: game_teams_enabled},
                {silent:true}
                );    
        
            // Save model
            var url = this.model.gameUrl();
            $(".save-game").text('Saving...');
            this.model.save(null, {
                url: url,
                success: function(model){
                    $(".save-game").text('Saved');        
                    setTimeout(function(){
                        $(".save-game").text('Save');    
                    }, 1000)
                }                
            });
        },
        startGameModal: function(event) {
          $('#start_game_modal').modal();
        },
        startGame: function(event) {
            $('#start_game_modal').modal('hide');
            var that = this;
            var url = this.model.gameUrl();
            $.post(url, function(){
                that.model.set('game_started', true);
            }).error(function(response){
                alert(response.responseText);
            });
        },
        endGameModal: function(event) {
          $('#end_game_modal').modal();
        },
        endGame: function(event) {
            $('#end_game_modal').modal('hide');
            var that = this;
            var url = this.model.gameUrl();

            this.model.destroy({
                url: url,
                success: function() {
                    if (!app.Running.Games.setArbitraryActiveGame()) {
                        Backbone.history.navigate('#logout', {
                            trigger: true
                        });
                        return;
                    }
                }
            });
        },
                
        render: function(){
            $('.modal-backdrop').remove();            
            var data = this.model.attributes;
            data.teams_enabled = data.game_properties.teams_enabled == 'true';
            this.$el.html(this.template(data))
            return this;
        }    
    })
})(jQuery);
    
//
// js/views/admin-user-view.js
// dmassassins.js
//
// Copyright (c) 2014 Matt Gerstman
// MIT License.
//
// displays user profile


var app = app || {
    Collections: {},
    Models: {},
    Views: {},
    Routers: {},
    Running: {},
    Session: {}
};

(function($) {
    'use strict';
    app.Views.AdminUserView = Backbone.View.extend({

        template: _.template($('#admin-user-template').html()),
        tagName:'div',
        initialize: function(model){
            this.model = model;
            this.listenTo(this.model, 'change', this.render)
            this.listenTo(this.model, 'fetch', this.render)
            this.listenTo(this.model, 'save', this.render)
        },
        render: function(extras){
            var data = this.model.attributes;
            for (var key in extras) {
                data[key] = extras[key]
            }
            this.$el.html(this.template(data))
            return this;
        }    
    })
})(jQuery);
    
//
// js/views/nav-game-view.js
// dmassassins.js
//
// Copyright (c) 2014 Matt Gerstman
// MIT License.
//
// handles the game dropdown in the nav


var app = app || {
    Collections: {},
    Models: {},
    Views: {},
    Routers: {},
    Running: {},
    Session: {}
};

(function($) {
    'use strict';
    app.Views.AdminUsersTeamsView = Backbone.View.extend({
    
        template: _.template($('#admin-users-teams-template').html()),
        tagName: 'ul',
        initialize: function() {
            this.collection = app.Running.Teams;
            this.listenTo(this.collection, 'fetch', this.render)
            this.listenTo(this.collection, 'change', this.render)
            this.listenTo(this.collection, 'reset', this.render)
            this.listenTo(this.collection, 'add', this.render)
            this.listenTo(this.collection, 'remove', this.render)
        },
        render: function() {
            var teamSort = function(team){
                return team.team_name;
            }
            
            var data = { teams: _.sortBy(this.collection.toJSON(), teamSort) };
                        
            var myRole = app.Running.User.getProperty('user_role');
            data.is_admin = AuthUtils.requiresAdmin(myRole);
            this.$el.html(this.template(data));
            return this;

        },

    })

})(jQuery);
//
// js/views/admin-users-view.js
// dmassassins.js
//
// Copyright (c) 2014 Matt Gerstman
// MIT License.
//
// displays user profile


var app = app || {
    Collections: {},
    Models: {},
    Views: {},
    Routers: {},
    Running: {},
    Session: {}
};

(function($) {
    'use strict';
    app.Views.AdminUsersView = Backbone.View.extend({


        template: _.template($('#admin-users-template').html()),
        tagName: 'div',

        // The DOM events specific to an item.
        events: {
            'click .ban-user': 'banUserModal',
            'click .kill-user': 'killUserModal',
            'click .revive-user': 'reviveUserModal',
            'click .ban-user-submit': 'banUser',
            'click .kill-user-submit': 'killUser',
            'click .revive-user-submit': 'reviveUser',
            'change select.user-team': 'selectChangeTeam',
            'change select.user-role': 'selectChangeRole',
            'click  a.team-name ': 'sortByTeam',
            'click .new-team-open': 'showNewTeam',
            'click .create-new-team': 'createNewTeam',
            'click .cancel-new-team': 'cancelNewTeam',
            'keyup .new-team-name': 'newTeamKeypress',
            'blur .new-team-form input': 'blurTeamForm',
            'click .edit-team': 'showEditTeamForm',
            'click .cancel-edit-team': 'cancelEditTeam',
            'click .save-edit-team': 'saveEditTeam',
            'click .delete-team': 'deleteTeamModal',
            'click .delete-team-submit': 'deleteTeam'
        },
        team: undefined,
        // constructor
        initialize: function() {
            var myRole = app.Running.User.getProperty('user_role');
            this.collection = app.Running.Users;
            this.userViews = [];
            this.teams_view = new app.Views.AdminUsersTeamsView();
            this.listenTo(this.collection, 'fetch', this.render);
            this.listenTo(this.collection, 'sync', this.render);
            this.listenTo(this.collection, 'change', this.render)
			this.listenTo(this.collection, 'remove', this.render)
        },
        banUserModal: function(event) {
            var user_name = $(event.currentTarget).data('user-name');
            var user_id = $(event.currentTarget).data('user-id');
            $('.ban-user-submit').data('user-id', user_id);
            $('#ban_user_modal .user-name').text(user_name)
            $('#ban_user_modal').modal();
        },
        killUserModal: function(event) {
            var user_name = $(event.currentTarget).data('user-name');
            var user_id = $(event.currentTarget).data('user-id');
            $('.kill-user-submit').data('user-id', user_id);
            $('#kill_user_modal .user-name').text(user_name)
            $('#kill_user_modal').modal();  
        },
        reviveUserModal: function(event) {
            var user_name = $(event.currentTarget).data('user-name');
            var user_id = $(event.currentTarget).data('user-id');
            $('.revive-user-submit').data('user-id', user_id);
            $('#revive_user_modal .user-name').text(user_name)
            $('#revive_user_modal').modal();  
        },
               
        banUser: function(event) {
        	var user_id = $(event.currentTarget).data('user-id');
	      	var user = this.collection.get(user_id);
	      	user.ban()
		  	$('#ban_user_modal').modal('hide')
        },
        killUser: function(event) {
        	var user_id = $(event.currentTarget).data('user-id');
	      	var user = this.collection.get(user_id);
	      	user.kill()
	      	$('#kill_user_modal').modal('hide')
        },

        reviveUser: function(event) {
        	var user_id = $(event.currentTarget).data('user-id');
	      	var user = this.collection.get(user_id);
	      	user.revive()
	      	$('#revive_user_modal').modal('hide')
        },

        selectChangeTeam: function(event){
            var user_id = $(event.currentTarget).data('user-id');
            var team_id = $(event.currentTarget).find('option:selected').val()
            var team_name = $(event.currentTarget).find('option:selected').text();
            this.addUserToTeam(user_id, team_id, team_name);
            
        },
        selectChangeRole: function(event){
            var user_id = $(event.currentTarget).data('user-id');
            var role_id = $(event.currentTarget).find('option:selected').val()
            return this.changeUserRole(user_id, role_id);
        },
        changeUserRole: function(user_id, role_id){
            // Sorry Taylor, a model for this one is overkill
            var game_id = app.Running.Games.getActiveGameId();
            var url = config.WEB_ROOT + 'game/' + game_id + '/user/' + user_id + '/role/';
            $.post(url, {role: role_id} ,function(){
                $('#role_saved_'+user_id).fadeIn(500, function(){ $(this).fadeOut(2000) })
            });
        },
        addUserToTeam: function(user_id, team_id, team_name, callback) {
            var that = this;
            var team = new app.Models.Team({user_id: user_id, team_id: team_id})
            var user = app.Running.Users.get(user_id);
            return team.save(null, {
                success: function(){
                    that.collection.get(user_id).setProperty('team', team_name);
                    $('#team_saved_'+user_id).fadeIn(500, function(){ $(this).fadeOut(2000) })
                    if (that.team !== undefined) {
                        $('#user_'+user_id).remove();
                    }
                }
            });                         
        },
        makeDraggable: function() {
            var that = this;            
            var startFunc = function(e, ui) {
                ui.helper.find('.user').remove();
                ui.helper.removeClass('user-grid');
                ui.helper.find('.drag-img').removeClass('hide');
                ui.helper.find('.drag-img').animate({
                    width: 50,
                    height: 50             
                }, 100);
            };
            
            this.$el.find('.user-grid').draggable({
                handle: '.thumbnail',
                connectWith: '#team_list li',
                tolerance: "pointer",
                helper: 'clone',
                forceHelperSize: true,
                zIndex:5000,
                start: startFunc,
                cursorAt: {left:40, top:25}
            })
        },
        makeDroppable: function() {
            var that = this;
            this.$el.find('#team_list li.team-droppable').droppable({
                hoverClass: 'drop-hover',
                tolerance: "pointer",
                drop: function(event, ui) {
                    var user_id = ui.helper.data('user-id');
                    var team_id = $(this).data('team-id');
                    var team_name = $(this).data('team-name');
                    that.addUserToTeam(user_id, team_id, team_name);
                }
            });
        },
        addUser: function(user, extras){
            extras.logged_in = false;
            if (user.get('user_id') == app.Session.get('user_id')) {
                extras.logged_in = true;
            }
            var userView = new app.Views.AdminUserView(user);
            this.userViews.push(userView);
            this.$el.find('.admin-users-body').append(userView.render(extras).el);
        },
        sortByTeam: function(event) {
            event.preventDefault();        
            this.team = $(event.currentTarget).data('team-name');
            this.team_id = $(event.currentTarget).data('team-id');
            if (this.team_id == 'SHOW_ALL') {
                this.team = undefined;
                this.team_id = 'all';
            }
                
                
            if (this.team_id == 'NO_TEAM') {
                this.team = "null";
                this.team_id = 'null'
            }
                
            
            this.render();            
        },
        showNewTeam: function(event) {            
            event.preventDefault();
            this.$el.find('.new-team-open').addClass('hide');
            this.$el.find('.new-team-form').removeClass('hide');
            this.$el.find('.new-team-form input').focus();
        },
        hideNewTeam: function() {             
            this.$el.find('.new-team-open').removeClass('hide');
            this.$el.find('.new-team-form').addClass('hide');
        },
        cancelNewTeam: function(event) {
            if (event)           
                event.preventDefault();
            this.hideNewTeam();
        },
        blurTeamForm: function() {
            var team_name = this.$el.find('.new-team input').val();
            if (!team_name) {                        
                this.hideNewTeam();    
            }
            
        },
        createNewTeam: function(event) {
            if (event)           
                event.preventDefault();
                
            var team_name = this.$el.find('.new-team input').val();
            if (!team_name) {                        
                return;
            }
            var game_id = app.Running.Games.getActiveGameId();
            var url = config.WEB_ROOT + 'game/' + game_id + '/team/';
            var that = this;
            $.post(url, {team_name:team_name}, function(team){
                app.Running.Teams.add(team);
                that.teams_view.render();
                that.selectActiveTeam();
                that.makeDroppable();
            });

        },
        newTeamKeypress: function(event) {
             if (event.keyCode == 27) {
                 this.hideNewTeam();                 
             }
             if (event.keyCode == 13) {
                 this.createNewTeam(event); 
             }
             
        },
        showEditTeamForm: function(event) {
            event.preventDefault();
            var team_id = $(event.currentTarget).data('team-id');
            $('#nav_team_'+team_id).find('.edit-team-form').removeClass('hide');
            $('#nav_team_'+team_id).find('.team-display').addClass('hide');

        },
        hideEditTeam: function(event) {
            var team_id = $(event.currentTarget).data('team-id');
            $('#nav_team_'+team_id).find('.team-display').removeClass('hide');
            $('#nav_team_'+team_id).find('.edit-team-form').addClass('hide');
        },
        cancelEditTeam: function(event) {
            if (event)           
                event.preventDefault();
            this.hideEditTeam(event);
        },
        saveEditTeam: function(event) {
            event.preventDefault();
            var team_id = $(event.currentTarget).data('team-id');
            var team = app.Running.Teams.get(team_id);
            var name = $('#nav_team_'+team_id).find('.edit-team-name').val()
            if (name == team.get('team_name'))
            {
                this.hideEditTeam(event);
                return;
            }
                
            var that = this;
            team.set('team_name', name);
            team.save();
        },
        deleteTeamModal: function(event) {
            event.preventDefault();
            var team_name = $(event.currentTarget).data('team-name');
            var team_id = $(event.currentTarget).data('team-id');
            $('.delete-team-submit').data('team-name', team_name);
            $('.delete-team-submit').data('team-id', team_id);
            $('#delete_team_modal .team-name').text(team_name)
            $('#delete_team_modal').modal();

        },
        deleteTeam: function(event) {
            var team_id = $(event.currentTarget).data('team-id');
            var team_name = $(event.currentTarget).data('team-name');
	      	var team = app.Running.Teams.get(team_id);
	      	var that = this;
	      	team.destroy({success:function(){
	      	    if (that.team == team_name) {
    	      	    that.team = 'null'
    	      	    that.team_id = 'null'
	      	    }
    	      	app.Running.Users.each(function(user){
        	      	if (user.getProperty('team') == team_name)
        	      	{
            	      	user.setProperty('team', 'null');
   
        	      	}
    	      	})
    	      	that.render();    	      	
            }});           
            $('#delete_team_modal').modal('hide');
        },
        selectActiveTeam: function() {
            this.$el.find('.active').removeClass('active');
            this.$el.find('#nav_team_'+this.team_id).addClass('active');

        },
        render: function() {			
			var data = {};
			var game = app.Running.Games.getActiveGame();
			var teams_enabled = false
			if (game)
			{
    			teams_enabled = game.areTeamsEnabled();
			}
            data.teams_enabled = teams_enabled;
			
            this.$el.html(this.template(data));
            
            this.teams_view.setElement(this.$('#team_list')).render();
                    
            while (this.userViews.length)
            {   
                var view = this.userViews.pop();
                view.remove();
            }

            var data = this.collection.models;
            if (this.team !== undefined)
            {
                var that = this;                
                data = _.filter(data, function(user){
                    return user.getProperty('team') == that.team;
                });                            
            }
            this.selectActiveTeam();

            var userSort = function(user) {
                return user.getProperty('first_name');
            }

            data = _.sortBy(data, userSort);        

            var myRole = app.Running.User.getProperty('user_role');
            var that = this;
            var extras = {
                teams: app.Running.Teams.toJSON(),
                roles: AuthUtils.getRolesMapFor(myRole, teams_enabled),
                is_admin: AuthUtils.requiresAdmin(myRole),
                teams_enabled: teams_enabled
            };    
                         
            _.each(data, function(user){
                that.addUser(user, extras);
            })

            if (teams_enabled)
            {
                this.makeDraggable();
                this.makeDroppable();                
            }
            this.trigger('render');
            return this;
        }
    })
})(jQuery);
//
// js/views/app-view.js
// dmassassins.js
//
// Copyright (c) 2014 Matt Gerstman
// MIT License.
//
// loads pages within the body of the app

var app = app || {
    Collections: {},
    Models: {},
    Views: {},
    Routers: {},
    Running: {},
    Session: {}
};

(function($) {
    'use strict';
    app.Views.AppView = Backbone.View.extend({
        el: '#app',
        // constructor
        initialize: function() {
            this.$body = $('#app_body');
        },
        // renders a page within the body of the app
        renderPage: function(page) {
            // Removes modal backdrop if we rapidly change pages
            $('.modal-backdrop').remove();
            this.$body.html(page.render().el);
        },
        setCurrentView: function(view) {
            if (app.Running.currentView)
                app.Running.currentView.remove();
            app.Running.currentView = view;

        }
    })
})(jQuery);
//
// js/views/leaderboard-view.js
// dmassassins.js
//
// Copyright (c) 2014 Matt Gerstman
// MIT License.
//
// shows the list of high scores


var app = app || {
    Collections: {},
    Models: {},
    Views: {},
    Routers: {},
    Running: {},
    Session: {}
};

(function($) {
    'use strict';
    app.Views.LeaderboardView = Backbone.View.extend({

        template: _.template($('#leaderboard-template').html()),
        tagName: 'div',

        // constructor
        initialize: function(params) {
            this.model = app.Running.LeaderboardModel;
            this.listenTo(this.model, 'change', this.render)
            this.listenTo(this.model, 'reset', this.render)
            this.listenTo(this.model, 'fetch', this.render)
        },
        // renderer
        render: function() {
            var data = this.model.attributes;
            this.$el.html(this.template(data));
            var numCols = 2;
            var teams_enabled = data.teams_enabled;
            if (teams_enabled)
                numCols = 3;
                
            var options = {
                paging: false,
                searching: false,
                info: false,
                order: [
                    [numCols - 1, 'desc'],
                    [numCols, 'desc']
                ]
            };

            this.$el.find('#user_leaderboard_table').dataTable(options);

            if (teams_enabled) {
                options.order = [
                    [4, 'desc']
                ]
                this.$el.find('#team_leaderboard_table').dataTable(options);
            }
            return this;
        }

    })

})(jQuery);
//
// js/views/login-view.js
// dmassassins.js
//
// Copyright (c) 2014 Matt Gerstman
// MIT License.
//
// shows the login screen


var app = app || {
    Collections: {},
    Models: {},
    Views: {},
    Routers: {},
    Running: {},
    Session: {}
};

(function($) {
    'use strict';
    app.Views.LoginView = Backbone.View.extend({
    
        template: _.template($('#login-template').html()),

        events: {
            'click .btn-facebook': 'login'
        },

        initialize: function() {
            this.model = app.Session;
        },
        // call the model login function
        login: function() {

            this.model.login()
        },
        // render the login page
        render: function() {
            this.$el.html(this.template(this.model.attributes));
            return this;
        },
    })

})(jQuery);
//
// js/views/nav-game-view.js
// dmassassins.js
//
// Copyright (c) 2014 Matt Gerstman
// MIT License.
//
// handles the game dropdown in the nav


var app = app || {
    Collections: {},
    Models: {},
    Views: {},
    Routers: {},
    Running: {},
    Session: {}
};

(function($) {
    'use strict';
    app.Views.NavGameView = Backbone.View.extend({
        template: _.template($('#nav-game-template').html()),
        el: '#games_dropdown',

        tagName: 'ul',

        events: {
            'click li a.switch_game': 'select'
        },
        // constructor, loads a user id so we can get their games from the model
        initialize: function() {
            this.collection = app.Running.Games;
            this.listenTo(this.collection, 'fetch', this.render)
            this.listenTo(this.collection, 'change', this.render)
            this.listenTo(this.collection, 'reset', this.render)
            this.listenTo(this.collection, 'add', this.render)
            this.listenTo(this.collection, 'remove', this.render)
            this.listenTo(this.collection, 'game-change', this.render);

        },
        handleJoin: function () {
            var availableGame = _.findWhere(this.collection.toJSON(), {member: false})
            if (availableGame === undefined)
            {
                this.hideJoin();
                return;
            }
            this.showJoin();
        },
        hideJoin: function () {
            this.$el.find('#nav_join_game').addClass('hide');
        },
        showJoin: function () {
            this.$el.find('#nav_join_game').removeClass('hide');
        },
        showCurrentGame: function() {
            var game_id = app.Running.Games.getActiveGameId();
            this.$el.find('#nav_' + game_id).removeClass('hide');
        },
        updateText: function() {

            $('.game_name').removeClass('hide');
            if (Backbone.history.fragment == 'join_game') {
                this.showCurrentGame();
                $('#games_header').text('Join Game');
                $('#games_header_short').text('Join Game');
                return this;
            }

            if (Backbone.history.fragment == 'create_game') {
                this.showCurrentGame();
                $('#games_header').text('Create Game');
                $('#games_header_short').text('Create Game');
                return this;
            }

            var game = this.collection.getActiveGame()
            if (!game)
            {
                return this;
            }
            var game_name = game.get('game_name');
            $('#games_header').text(game_name);
            var max = 9;
            if (game_name.length > max) {
                game_name = game_name.substr(0, max - 3) + '...';
            }
            $('#games_header_short').text(game_name);

            var game_id = app.Running.Games.getActiveGameId();
            this.$el.find('#nav_' + game_id).addClass('hide');
        },
        // loads the items into the dropdown and changes the dropdown title to the current game
        render: function() {
            this.$el.html(this.template({
                games: _.where(this.collection.toJSON(), {
                    member: true
                })
            }));
            this.handleJoin();
            this.updateText();
            return this;

        },
        // select a game from the dropdown
        select: function(event) {
            var game_id = $(event.target).attr('game_id');
            app.Running.Games.setActiveGame(game_id)

        }
    })

})(jQuery);
//
// js/views/nav-view.js
// dmassassins.js
//
// Copyright (c) 2014 Matt Gerstman
// MIT License.
//
// handles the nav bar at the top


var app = app || {
    Collections: {},
    Models: {},
    Views: {},
    Routers: {},
    Running: {},
    Session: {}
};

(function($) {
    'use strict';
    app.Views.NavView = Backbone.View.extend({


        template: _.template($('#nav-template').html()),
        el: '#nav_body',

        tagName: 'nav',

        events: {
            'click li a': 'select'
        },

        // constructor
        initialize: function() {
            this.NavGameView = app.Running.NavGameView;
            this.listenTo(app.Running.TargetModel, 'fetch', this.handleTarget)
            this.listenTo(app.Running.TargetModel, 'change', this.handleTarget)
            this.listenTo(app.Running.User, 'fetch', this.render)
            this.listenTo(app.Running.User, 'change', this.render)
            this.listenTo(app.Running.Games, 'game-change', this.render)
        },

        // if we don't have a target hide that view
        render: function() {
            var role = app.Running.User.getProperty('user_role');  
            var data = {};
            data.is_captain = AuthUtils.requiresCaptain(role);
            data.is_admin = AuthUtils.requiresAdmin(role);
            
            this.$el.html(this.template(data));
            this.handleTarget();
            
            var selectedElem = this.$el.find('#nav_' + Backbone.history.fragment);
            this.highlight(selectedElem);
            
            if (app.Running.NavGameView)
                app.Running.NavGameView.setElement(this.$('#games_dropdown')).render();
            return this;
        },

        // select an item on the nav bar
        select: function(event) {
            var target = event.currentTarget;
            if ($(target).hasClass('disabled') || $(target).hasClass('dropdown-toggle')) {
                event.preventDefault();
                return;
            }
            $('.navbar-collapse.in').collapse('hide');
            this.highlight(target)

        },

        // highlight an item on the nav bar and unhighlight the rest of them
        highlight: function(elem) {
            if ($(elem).hasClass('dropdown_parent')) {
                return;
            }

            if ($(elem).hasClass('dropdown_item')) {
                var dropdown = $(elem).attr('dropdown');
                var parent = '#' + dropdown + '_parent';
                elem = parent;
            }
            $('.active').removeClass('active');
            $(elem).addClass('active');
        },
        handleAdmin: function() {
            var role = app.Running.User.getProperty('user_role');  
            var allowed = AuthUtils.requiresCaptain(role);
            if (allowed) {
                $('#admin_parent').removeClass('hide');
                return;
            }
            $('#admin_parent').addClass('hide');
            return;
            
        },
        handleTarget: function() {
            var game = app.Running.Games.getActiveGame();
            if (!game)
            {
                this.disableTarget();
                return;
            }
            if (!game.get('game_started'))
            {
                this.disableTarget();
                return;
            }
            
            if (!app.Running.TargetModel.get('user_id'))
            {
                this.disableTarget();
                return;
            }
            
            this.enableTarget();
            return;
        },

        // hides the target nav item
        enableTarget: function() {
            this.$el.find('#nav_target').removeClass('disabled');
            this.$el.find('#nav_target a').removeClass('disabled');
        },

        // shows the target nav item
        disableTarget: function() {
            this.$el.find('#nav_target').addClass('disabled');
            this.$el.find('#nav_target a').addClass('disabled');
        }

    })

})(jQuery);
//
// js/views/profile-view.js
// dmassassins.js
//
// Copyright (c) 2014 Matt Gerstman
// MIT License.
//
// displays user profile


var app = app || {
    Collections: {},
    Models: {},
    Views: {},
    Routers: {},
    Running: {},
    Session: {}
};

(function($) {
    'use strict';
    app.Views.ProfileView = Backbone.View.extend({


        template: _.template($('#profile-template').html()),
        tagName: 'div',

        // The DOM events specific to an item.
        events: {
            'click .thumbnail': 'showFullImage',
            'click #quit': 'showQuitModal',
            'click #quit_game_confirm': 'quitGame',
            'click #email_settings': 'showEmailModal',
            'click #email_settings_save': 'saveEmailSettings',
            'keyup #email': 'emailEnter' 
        },

        // load profile picture in modal window
        showFullImage: function() {
            $('#photoModal').modal()
        },
        // load quit confirm modal
        showQuitModal: function() {
            var templateVars = {
                quit_game_name: app.Running.Games.getActiveGame().get('game_name')
            }
            var template = _.template($('#quit-modal-template').html())
            var html = template(templateVars);
            $('#quit_modal_wrapper').html(html);
            $('#quit_modal').modal();
        },
        quitGame: function() {
            var secret = this.$el.find('#quit_secret').val();
            this.model.quit(secret);

        },
        showEmailModal: function(){
            $('#email_modal').modal();
        },
        emailEnter: function(event) {
            if (event.which == 13) {
              this.saveEmailSettings();
            }  
        },
        saveEmailSettings: function(){
            var email = $('#email').val();
            var allow_email = $('#allow_email').is(':checked') ? 'true' : 'false';
            this.model.saveEmailSettings(email, allow_email);
            $('#email_modal').modal('hide')
        },
        // constructor
        initialize: function(params) {
            this.model = app.Running.User;
            this.listenTo(this.model, 'change', this.render)
            this.listenTo(this.model, 'fetch', this.render)
            this.listenTo(this.model, 'destroy', this.destroyCallback)
            this.listenTo(this.model, 'set', this.render)

        },
        destroyCallback: function() {
            $('#quit_modal').modal('hide')            
        },
        render: function() {
            $('.modal-backdrop').remove();
            var data = this.model.attributes;
            data.teams_enabled = false;
            var game = app.Running.Games.getActiveGame();
            if (game) {
                data.teams_enabled = game.areTeamsEnabled();    
            }
            data.allow_email = data.properties.allow_email == 'true';
            
            var role = app.Running.User.getProperty('user_role');  
            var allow_quit = !AuthUtils.requiresCaptain(role);
            data.allow_quit = allow_quit;

            this.$el.html(this.template(data));
            return this;
        }
    })
})(jQuery);
//
// js/views/profile-view.js
// dmassassins.js
//
// Copyright (c) 2014 Matt Gerstman
// MIT License.
//
// displays rules for a game


var app = app || {
    Collections: {},
    Models: {},
    Views: {},
    Routers: {},
    Running: {},
    Session: {}
};

(function($) {
    'use strict';
    app.Views.RulesView = Backbone.View.extend({


        template: _.template($('#rules-template').html()),
        tagName: 'div',

        initialize: function(params) {
            this.model = app.Running.RulesModel;

            this.listenTo(this.model, 'change', this.render)
            this.listenTo(this.model, 'fetch', this.render)
            this.listenTo(this.model, 'set', this.render)
        },

        render: function() {
            var data = this.model.attributes;
            data.rules = marked(data.rules);
            this.$el.html(this.template(data));
            return this;
        }

    })

})(jQuery);
//
// js/views/select-game-view.js
// dmassassins.js
//
// Copyright (c) 2014 Matt Gerstman
// MIT License.
//
// handles game selection


var app = app || {
    Collections: {},
    Models: {},
    Views: {},
    Routers: {},
    Running: {},
    Session: {}
};

(function($) {
    'use strict';
    app.Views.SelectGameView = Backbone.View.extend({


        template: _.template($('#select-game-template').html()),
        tagName: 'div',
        events: {
            'click .show-create-game': 'showCreateGame',
            'click .show-join-game': 'showJoinGame',
            'click .create-game-submit': 'createGame',
            'click .join-game-submit': 'joinGame',
            'click .create-or-join-back': 'goBack',
            'click #create_game_need_password': 'togglePassword',
            'change #games': 'checkFields'

        },
        // previous page, may depricate
        loaded_from: 'login',
        // constructor
        initialize: function() {
            this.collection = app.Running.Games;
            this.listenTo(this.collection, 'reset', this.render)
            this.listenTo(this.collection, 'fetch', this.render)
            this.listenTo(app.Running.User, 'join-error-password', this.badPassword)
        },
        // shows the create game subview
        showCreateGame: function() {
            $('.logo').addClass('hide');
            $('#create-or-join').addClass('hide');
            $('#create-game').addClass('select-game-active');
            $('#create-game').removeClass('hide');
        },
        // shows the join game subview
        showJoinGame: function() {
            $('.logo').addClass('hide');
            $('#create-or-join').addClass('hide');
            $('#join-game').addClass('select-game-active');
            $('#join-game').removeClass('hide');
        },
        // cancels the game creation/selection
        goBack: function() {
            if (!!app.Running.Games.getActiveGameId()) {
                app.Running.Router.back()
                return;
            }
            $('.select-game-active').addClass('hide').removeClass('select-game-active');
            $('#create-or-join').removeClass('hide');
            $('.logo').removeClass('hide');
        },
        // show the create game s ubview
        createGame: function(event) {
            event.preventDefault();
            var name = $('#create_game_name').val();
            var password = null;
            if ($('#create_game_need_password').is(':checked')) {
                password = $('#create_game_password').val();

            }
            var that = this;
            this.collection.create({
                game_name: name,
                game_password: password
            }, {
                success: function(game) {
                    that.finish(game);
                }
            });
        },
        // loads the join game later view
        loadJoinGame: function(user_id) {
            var that = this;
            that.showJoinGame();
            this.collection.fetch({
                wait: true,
                success: function() {
                    that.showJoinGame();
                }
            });

        },
        // posts to the join game model
        joinGame: function(event) {
            event.preventDefault();            
            var selected = this.$el.find('#games option:selected');
            var game_id = selected.val();
            var need_password = selected.attr('game_has_password') == 'true';
            var password = need_password ? $('#join_game_password').val() : '';
            
            var teams_enabled = this.$el.find('#join_game_team').attr('disabled') != 'disabled';
            var team_id = this.$el.find('#join_game_team option:selected').val();
            
            if (teams_enabled && !team_id) {
                this.badTeam();
                return;
            }            
            app.Running.Games.joinGame(game_id, password, team_id);
        },
        badPassword: function(){
            $('#join_password_block').addClass('has-error');
            $('label[for=join_game_password]').text('Invalid Password:');
        },
        badTeam: function(){
            $('#join_team_block').addClass('has-error');
            $('label[for=join_game_team]').text('Must Select A Team:');
        },
        fixFields: function(){
          this.$el.find('.has-error').removeClass('has-error');
          this.$el.find('label[for=join_game_password]').text('Password:');
          this.$el.find('label[for=join_game_team]').text('Team:');
        },
        // finish up and navigate to your profile
        finish: function(game) {
            app.Running.Games.setActiveGame(game.get('game_id'));
            Backbone.history.navigate('my_profile', {
                trigger: true
            });
        },
        // toggles the password entry field on create game
        togglePassword: function(e) {
            $('#create_game_password').attr('disabled', !e.target.checked);
        },
        noTeams: function(){
            var teamField = this.$el.find('#join_game_team');
            teamField.attr('disabled', true);            
            teamField.find('#team_placeholder').text('This Game Doesn\'t Have Teams');
        },
        // toggles the password entry field and team entry fieldon join game
        checkFields: function() {
            // Password field setup
            this.fixFields();
            var selected = this.$el.find('#games option:selected');
            var need_password = selected.attr('game_has_password') == 'true';
            var passwordField = this.$el.find('#join_game_password');
            var passwordPlaceholder = need_password ? '' : 'This Game Has No Password';
            passwordField.attr('disabled', !need_password);
            passwordField.val(passwordPlaceholder);
            
            var game_id = selected.val();
            var game = app.Running.Games.get(game_id);
            if (!game) {
                return;
            }

            var teamField = this.$el.find('#join_game_team');
            teamField.find('#team_placeholder').text('Loading..');

            var that = this;
            var url = config.WEB_ROOT + 'game/' + game_id + '/team/';
            app.Running.Teams.fetch({
                url:url,
                success:function(teams){        
                    teamField.attr('disabled', false);
                    var teamOptionsTemplate = _.template($('#select-game-team-option').html());
                    var teamOptionsHTML = teamOptionsTemplate({teams: app.Running.Teams.toJSON()});
                    teamField.html(teamOptionsHTML);
                    if (!app.Running.Teams.length){
                        that.noTeams();
                    }
                },
                error:function(){
                   that.noTeams()
                }
            });
            
        },
        render: function() {
            this.$el.html(this.template({
                games: _.where(this.collection.toJSON(), {
                    member: false
                })
            }));
            this.checkFields();
            return this;
        }

    })

})(jQuery);
//
// js/views/target-view.js
// dmassassins.js
//
// Copyright (c) 2014 Matt Gerstman
// MIT License.
//
// target view


var app = app || {
    Collections: {},
    Models: {},
    Views: {},
    Routers: {},
    Running: {},
    Session: {}
};

(function($) {
    'use strict';
    app.Views.TargetView = Backbone.View.extend({


        template: _.template($('#target-template').html()),
        tagName: 'div',

        // The DOM events specific to an item.
        events: {
            'click .thumbnail': 'showFullImage',
            'click #kill': 'kill'
        },
        // loads picture in a modal window
        showFullImage: function() {
            $('#photoModal').modal()
        },
        // constructor
        initialize: function() {
            this.model = app.Running.TargetModel;
            this.listenTo(this.model, 'change', this.render)
            this.listenTo(this.model, 'fetch', this.render)
            this.listenTo(this.model, 'set', this.render)
        },
        // kills your target
        kill: function() {
            var secret = this.$el.find('#secret').val();
            var view = this;
            this.model.destroy({
                headers: {
                    'X-DMAssassins-Secret': secret
                },
                success: function() {
                    view.model.fetch();
                }
            })
        },
        render: function() {
            var data = this.model.attributes;
            data.teams_enabled = app.Running.Games.getActiveGame().areTeamsEnabled();
            this.$el.html(this.template(data));
            return this;
        }
    })

})(jQuery);
// Base Router
// Provided From https://github.com/DanialK/advanced-security-in-backbone

var app = app || {
    Collections: {},
    Models: {},
    Views: {},
    Routers: {},
    Running: {},
    Session: {}
};

(function(){

	app.Routers.BaseRouter = Backbone.Router.extend({
		before: function(){},
		after: function(){},
		route : function(route, name, callback){
			if (!_.isRegExp(route)) route = this._routeToRegExp(route);
			if (_.isFunction(name)) {
				callback = name;
				name = '';
		 	}
		  	if (!callback) callback = this[name];

		  	var router = this;

		  	Backbone.history.route(route, function(fragment) {
		   		var args = router._extractParameters(route, fragment);
		   		var next = function(){
			    	callback && callback.apply(router, args);
				    router.trigger.apply(router, ['route:' + name].concat(args));
				    router.trigger('route', name, args);
				    Backbone.history.trigger('route', router, name, args);
				    router.after.apply(router, args);		
		   		}
		   		router.before.apply(router, [args, next]);
		  	});
			return this;
		}
	});
})()

//
// js/routers/router.js
// dmassassins.js
//
// Copyright (c) 2014 Matt Gerstman
// MIT License.
//
// Handles all the URL *Magic*

var app = app || {
    Collections: {},
    Models: {},
    Views: {},
    Routers: {},
    Running: {},
    Session: {}
};

(function() {

    app.Routers.Router = app.Routers.BaseRouter.extend({

        // Sometime we wanna go back and tis is the only wat to do that.
        history: [],
        // All the routes
        routes: {
            '': 'my_profile',
            'login': 'login',
            'logout': 'logout',
            'target': 'target',
            'my_profile': 'my_profile',
            'multigame': 'multigame',
            'leaderboard': 'leaderboard',
            'create_game': 'create_game',
            'join_game': 'join_game',
            'rules': 'rules',
            'users': 'users',
            'edit_rules': 'edit_rules',
            'switch_game': 'switch_game',
            'game_settings': 'game_settings',

        },
        // routes that require we have a game that has been started
        requiresTarget: ['#target'],

        // routes that require just authentication
        requiresJustAuth: ['#multigame', ''],

        // routes that require we have a game and we're authenticated
        requiresGameAndAuth: ['#my_profile', '#join_game', '#leaderboard', '#rules'],

        // routes that require the user is at least a team captain
        requiresCaptain: ['#users'],
        
        // routes that require is at least a game admin
        requiresAdmin: ['#edit_rules', '#game_settings', '#plot_twists'],        

        // routes that should hide the nav bar
        noNav: ['login', 'multigame'],

        // routes that a logged in user can't access
        preventAccessWhenAuth: ['#login'],

        // place to redirect users for requiresGameAndAuth
        redirectWithoutGame: '#multigame',

        // place to redirect logged in users who don't have a started game
        redirectWithoutGameStarted: 'my_profile',

        // place to redirect users who aren't logged in
        redirectWithoutAuth: '#login',

        before: function(params, next) {

            // is the user authenticated
            var isAuth = app.Session.get('authenticated');
            var path = Backbone.history.location.hash;

            // do we need a game and authentication
            var needGameAndAuth = _.contains(this.requiresGameAndAuth, path);

            // do we need authentication
            var needAuth = _.contains(this.requiresJustAuth, path);

            // should we prevent this page if authorized
            var cancelAccess = _.contains(this.preventAccessWhenAuth, path);

            // does this route need a running game
            var needTarget = _.contains(this.requiresTarget, path);

            // is there a game
            var hasGame = app.Session.get('has_game') == "true";

            // is the game started
            var hasTarget = !!app.Running.TargetModel.get('user_id') && app.Running.Games.getActiveGame().get('game_started');

            // do we need to be a captain
            var needCaptain = _.contains(this.requiresCaptain, path);

            // do we need to be an admin
            var needAdmin = _.contains(this.requiresAdmin, path);
            
            // The active user's role in the current game
            var userRole = app.Running.User.getProperty('user_role');
            
            // is the user a captain
            var isCaptain = AuthUtils.requiresCaptain(userRole);

            // is the user an admin
            var isAdmin = AuthUtils.requiresAdmin(userRole);



            /*
			Variables I use when shit's not routing properly */
            /*/
			console.log('path:', path);
			console.log('needGameAndAuth: ', needGameAndAuth);
			console.log('hasGame: ', hasGame);
			console.log('hasTarget: ', hasTarget);
			console.log('needAuth: ', needAuth);
			console.log('isAuth: ', isAuth);
			console.log('cancelAccess: ', cancelAccess);
            console.log('needCaptain: ', needCaptain);
            console.log('isCaptain: ', isCaptain);
            console.log('needAdmin: ', needAdmin);
            console.log('isAdmin: ', isAdmin);

/**/

            // Do we need authentication
            if ((needAuth || needGameAndAuth) && !isAuth) {
                app.Session.set('redirect_from', path);
                Backbone.history.navigate(this.redirectWithoutAuth, {
                    trigger: true
                });
            }
            // do we need authentication and a game
            else if (needGameAndAuth && !hasGame) {
                Backbone.history.navigate(this.redirectWithoutGame, {
                    trigger: true
                });
            }
            // do we need a game and is it started
            else if (needTarget && !hasTarget) {
                Backbone.history.navigate(this.redirectWithoutTarget, {
                    trigger: true
                });
            }
            // are they logged in and trying to hit the login page
            else if (isAuth && cancelAccess) {
                Backbone.history.navigate('', {
                    trigger: true
                });
            }
            // do they need to be a captain and are they?
            else if (needCaptain && !isCaptain) {
                Backbone.history.navigate('', {
                    trigger: true
                });
            }
            // do they need to be an admin and are they?
            else if (needAdmin && !isAdmin) {
                Backbone.history.navigate('', {
                    trigger: true
                });
            // nothing is wrong! let them pass.	
            } else {
                //No problem handle the route
                return next();
            }
        },
        // called after we're done routing, unused but build into the baserouter so we're leaving it
        after: function() {
            this.history.push(Backbone.history.fragment);
        },
        // go to the previous
        back: function() {
            var path = this.history.pop();
            Backbone.history.navigate(path, {
                    trigger: true
                });
        },
        // login route
        login: function() {
            var view = new app.Views.LoginView();
            app.Running.AppView.setCurrentView(view);
            this.render();
        },
        // logout route
        logout: function() {
            app.Session.clear()
            this.navigate('login', true)
        },
        // game selection route
        multigame: function() {
            var view = new app.Views.SelectGameView();
            app.Running.AppView.setCurrentView(view);
            app.Running.currentView.collection.fetch();
            this.render();
        },
        // target route
        target: function() {
            var view = new app.Views.TargetView();
            app.Running.AppView.setCurrentView(view);
            this.render();
        },
        // create a new game route
        create_game: function() {
            var view = new app.Views.SelectGameView();
            app.Running.AppView.setCurrentView(view);
            this.render();
            app.Running.currentView.showCreateGame();
        },
        // join a new game route
        join_game: function() {
            var view = new app.Views.SelectGameView();
            app.Running.AppView.setCurrentView(view);
            this.render();
            app.Running.currentView.loadJoinGame(app.Session.get('user_id'));
        },
        // profile route
        my_profile: function() {
            var view = new app.Views.ProfileView();
            app.Running.AppView.setCurrentView(view);
            this.render();
        },
        // leaderboard route
        leaderboard: function() {
            var view = new app.Views.LeaderboardView();
            app.Running.AppView.setCurrentView(view);
            app.Running.currentView.model.fetch();
            this.render();
        },
        // rules route
        rules: function() {
            var view = new app.Views.RulesView();
            app.Running.AppView.setCurrentView(view);
            this.render();
        },
        users: function() {            
            var view = new app.Views.AdminUsersView();
            app.Running.AppView.setCurrentView(view);            
            app.Running.currentView.collection.reset();
            app.Running.currentView.collection.fetch();
            app.Running.Teams.fetch();             
            this.render();
        },
        edit_rules: function() {
            var view = new app.Views.AdminEditRulesView();
            app.Running.AppView.setCurrentView(view);
            app.Running.currentView.model.fetch();
            this.render(); 
        },
        game_settings: function() {
            var view = new app.Views.AdminGameSettingsView();
            app.Running.AppView.setCurrentView(view);
            this.render();             
        },
        preventSwitchGameBack: ['join_game', 'create_game'],
        switch_game: function() {
            var lastFragment = this.history[this.history.length - 1];
            if (lastFragment === undefined || _.contains(this.preventSwitchGameBack, lastFragment)) {             
                Backbone.history.navigate('my_profile', {
                    trigger: true
                });
                return;
            }

            this.back();
        },
        // render function, also determines weather or not to render the nav
        render: function() {
            var fragment = Backbone.history.fragment;
            // if it's a view with a nav and we don't have one, make one
            if ((this.noNav.indexOf(Backbone.history.fragment) == -1) && (fragment != 'login') && (!app.Running.NavView)) {
                if (!app.Running.NavGameView) {
                    app.Running.NavGameView = new app.Views.NavGameView();
                    app.Running.NavGameView.render();
                }
                app.Running.NavView = new app.Views.NavView();
                app.Running.NavView = app.Running.NavView.render();
            }
            // if it explicitely shouldn't have a nav and we have one kill it
            else if ((this.noNav.indexOf(Backbone.history.fragment) != -1) && (app.Running.NavView)) {
                app.Running.NavView.$el.html('');
                app.Running.NavView = null;
                app.Running.navGameView = null;
            }
            // if we have a nav and highlight the nav item
            if ((app.Running.NavView) && (this.noNav.indexOf(Backbone.history.fragment) == -1)) {
                if (fragment === '')
                    fragment = 'my_profile';

                app.Running.NavView.highlight('#nav_' + fragment)
                app.Running.NavView.handleTarget();
                app.Running.NavGameView.updateText();
            }

            // render our page within the app
            app.Running.AppView.renderPage(app.Running.currentView)
        }
    })

})()
//
// app.js
// dmassassins.js
//
// Copyright (c) 2014 Matt Gerstman
// MIT License.
//

var app = app || {
    Collections: {},
    Models: {},
    Views: {},
    Routers: {},
    Running: {},
    Session: {}
};

// Instantiates all of the running models, routers, and session

$(function() {
    'use strict';

    Raven.config(config.SENTRY_DSN, {
    }).install();

    app.Running.AppView = new app.Views.AppView();
    app.Running.AppView.render();

    app.Session = new app.Models.Session();
    app.Session.setAuthHeader();

    app.Running.Games = new app.Collections.Games();
    app.Running.User = new app.Models.User()
    app.Running.TargetModel = new app.Models.Target()
    app.Running.LeaderboardModel = new app.Models.Leaderboard();
    app.Running.RulesModel = new app.Models.Rules();

    app.Running.Users = new app.Collections.Users();
    app.Running.Teams = new app.Collections.Teams();

    app.Running.User.listenTo(app.Running.Games, 'game-change', app.Running.User.fetch);
    app.Running.TargetModel.listenTo(app.Running.Games, 'game-change', app.Running.TargetModel.fetch);
    app.Running.LeaderboardModel.listenTo(app.Running.Games, 'game-change', app.Running.LeaderboardModel.fetch);
    app.Running.RulesModel.listenTo(app.Running.Games, 'game-change', app.Running.RulesModel.fetch);

    app.Running.User.listenTo(app.Running.User, 'fetch', app.Running.User.checkAccess);
    app.Running.User.listenTo(app.Running.User, 'change', app.Running.User.checkAccess);

    app.Running.Router = new app.Routers.Router();
    Backbone.history.start();

});
var config = {
    ENV: "DEV",
	LOCAL_ROOT : '/Users/Matthew/Dropbox/Go/src/GitHub.com/mattgerstman/DMAssassins/webapp/',
	WEB_ROOT : 'http://assassins.com:8000/',
	APP_ID : 649985815097586,
    SENTRY_DSN: 'https://b4ab4d3478f440e0bf658313a71b9847@app.getsentry.com/30968'
};
