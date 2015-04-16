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

        // Sometime we wanna go back and this is the only way to do that.
        history: [],
        // All the routes
        routes: {
            '': config.DEFAULT_VIEW,
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
            'targets': 'targets',
            'support': 'support'

        },
        // routes that require we have a game that has been started
        requiresTarget: ['', '#target'],

        // routes that require just authentication
        requiresJustAuth: ['#multigame'],

        // routes that require we have a game and we're authenticated
        requiresGameAndAuth: ['#my_profile', '#join_game', '#leaderboard', '#rules'],

        // routes that require the user is at least a team captain
        requiresCaptain: ['#users'],

        // routes that require is at least a game admin
        requiresAdmin: ['#edit_rules', '#game_settings', '#plot_twists', '#email_users'],

        // routes that require is a super admin
        requiresSuperAdmin: ['#targets'],

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

        // place to redirect to when a user doesn't have a target
        redirectWithoutTarget: '#my_profile',

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
            var sessionHasGame = app.Session.get('has_game');
            var hasGame = (sessionHasGame == "true") || (sessionHasGame === true);

            // do we need to be a captain
            var needCaptain = _.contains(this.requiresCaptain, path);

            // do we need to be an admin
            var needAdmin = _.contains(this.requiresAdmin, path);

            // do we need to be an admin
            var needSuperAdmin = _.contains(this.requiresSuperAdmin, path);

            // The active user's role in the current game
            var userRole = app.Running.User.getRole();

            // is the user a captain
            var isCaptain = AuthUtils.requiresCaptain(userRole);

            // is the user an admin
            var isAdmin = AuthUtils.requiresAdmin(userRole);

            // is the user a super admin
            var isSuperAdmin = AuthUtils.requiresSuperAdmin(userRole);

            // is the game started
            var gameStarted = app.Running.Games.getActiveGame() && app.Running.Games.getActiveGame().get('game_started');

            // is the game started
            var hasTarget = !!app.Running.TargetModel.get('user_id') && gameStarted && !isAdmin;

            /*
      Variables I use when shit's not routing properly */
            /*/
            console.log('userRole:', userRole);
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
            // do we need a target and do we have one
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
            } else if (needSuperAdmin && !isSuperAdmin) {
                Backbone.history.navigate('', {
                    trigger: true
                });
            }
            else {
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
            this.history.pop();
            var path = this.history.pop();
            console.log(path);
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
            app.Session.clear();
            this.navigate('login', true);
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
            app.Running.Teams.fetch();
            var view = new app.Views.AdminUsersView();
            app.Running.AppView.setCurrentView(view);
            app.Running.currentView.collection.fetch({reset: true});
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
        email_users: function() {
            var view = new app.Views.AdminEmailUsersView();
            app.Running.UserEmails.fetch();
            app.Running.AppView.setCurrentView(view);
            this.render();
        },
        targets: function() {
            var view = new app.Views.SuperAdminTargetsView();
            app.Running.AppView.setCurrentView(view);
            this.render();
            app.Running.currentView.model.fetch({reset: true});
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
        support: function() {
            this.navigate('', true);
            $('.js-support').click();
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
                    fragment = config.DEFAULT_VIEW;

                fragment = fragment.replace('_', '-');
                app.Running.NavView.highlight('#js-nav-' + fragment);
                app.Running.NavView.handleTarget();
                app.Running.NavGameView.updateText();
            }

            // render our page within the app
            app.Running.AppView.renderPage(app.Running.currentView);
        }
    });
})();
