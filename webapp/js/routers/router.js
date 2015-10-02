//
// js/routers/router.js
// dmassassins.js
//
// Copyright (c) 2014 Matt Gerstman
// MIT License.
//
// Handles all the URL *Magic*

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
            'my-profile': 'my_profile',
            'multigame': 'multigame',
            'leaderboard': 'leaderboard',
            'create-game': 'create_game',
            'join-game': 'join_game',
            'rules': 'rules',
            'users': 'users',
            'edit-rules': 'edit_rules',
            'switch-game': 'switch_game',
            'game-settings': 'game_settings',
            'targets': 'targets',
            'support': 'support',

            // old links that we want to alert about if an error occurs
            'my_profile': 'oldLink',
            'create_game': 'oldLink',
            'join_game': 'oldLink',
            'edit_roles': 'oldLink',
            'switch_game': 'oldLink',
            'game_settings': 'oldLink',

        },
        // routes that require we have a game that has been started
        requiresTarget: ['', 'target'],

        // routes that require just authentication
        requiresJustAuth: ['multigame'],

        // routes that require we have a game and we're authenticated
        requiresGameAndAuth: ['my-profile', 'leaderboard', 'rules'],

        // routes that require the user is at least a team captain
        requiresCaptain: ['users'],

        // routes that require is at least a game admin
        requiresAdmin: ['edit-rules', 'game-settings', 'plot-twists', 'email-users'],

        // routes that require is a super admin
        requiresSuperAdmin: ['targets'],

        // routes that should hide the nav bar
        noNav: ['login', 'multigame'],

        // routes that a logged in user can't access
        preventAccessWhenAuth: ['login'],

        // place to redirect users for requiresGameAndAuth
        redirectWithoutGame: 'multigame',

        // place to redirect logged in users who don't have a started game
        redirectWithoutGameStarted: 'my-profile',

        // place to redirect users who aren't logged in
        redirectWithoutAuth: 'login',

        // place to redirect to when a user doesn't have a target
        redirectWithoutTarget: 'my-profile',

        before: function(params, next) {

            // is the user authenticated
            var isAuth = app.Session.get('authenticated') === true;
            var path = Backbone.history.location.pathname;
            if (path.indexOf('/') === 0){
                path = path.substring(1);
            }

            // do we need a game and authentication
            var needGameAndAuth = _.contains(this.requiresGameAndAuth, path);

            // do we need authentication
            var needAuth = _.contains(this.requiresJustAuth, path);

            // should we prevent this page if authorized
            var cancelAccess = _.contains(this.preventAccessWhenAuth, path);

            // does this route need a running game
            var needTarget = _.contains(this.requiresTarget, path);

            // is there a game
            var hasGame = app.Session.get('has_game') === true;

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
            var gameStarted = app.Running.Games.hasActiveGameStarted();
            // is the game started
            var hasTarget = !!app.Running.TargetModel.get('user_id') && gameStarted && !isAdmin;

            /*
            Variables I use when shits not routing properly */
            /*/
            console.log('path:', path);
            console.log('gameStarted:', gameStarted);
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
                // No problem handle the route
                if (needCaptain) {
                    return app.Running.Async.requiresCaptain(next);
                }
                if (needAdmin) {
                    return app.Running.Async.requiresAdmin(next);
                }
                if (needSuperAdmin) {
                    return app.Running.Async.requiresSuperAdmin(next);
                }
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
            app.Session.logout();
            this.navigate('login', true);
        },
        // game selection route
        multigame: function() {
            var view = new app.Views.MultiGameView();
            app.Running.AppView.setCurrentView(view);
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
            var view = new app.Views.CreateGameView();
            app.Running.AppView.setCurrentView(view);
            this.render();
        },
        // join a new game route
        join_game: function() {
            var view = new app.Views.JoinGameView();
            app.Running.AppView.setCurrentView(view);
            this.render();

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
            // app.Running.currentView.collection.fetch({ reset: true });
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
        preventSwitchGameBack: ['join-game', 'create-game'],
        switch_game: function() {
            var lastFragment = this.history[this.history.length - 1];
            if (lastFragment === undefined || _.contains(this.preventSwitchGameBack, lastFragment)) {
                Backbone.history.navigate('my-profile', {
                    trigger: true
                });
                return;
            }

            this.back();
        },
        support: function() {
            this.navigate('', true);
            var supportView = new app.Views.SupportView();
            supportView.render();
        },
        oldLink: function() {
            console.log(arguments);
            alert("Invalid link");
        },
        // render function, also determines weather or not to render the nav
        render: function() {
            var fragment = Backbone.history.fragment;

            var hasActiveGame = app.Running.Games.getActiveGame() !== null;
            var recoveringSession = app.Session.get('has_game');

            // if we have a game, create a navbar
            if ((hasActiveGame || recoveringSession) && (fragment !== 'login') && (!app.Running.NavView)) {
                app.Running.NavView = new app.Views.NavView();
                app.Running.NavView = app.Running.NavView.render();
            }
            // if dont have a game, but we have a navbar delete it
            else if ((!hasActiveGame) && (app.Running.NavView)) {
                app.Running.NavView.$el.html('');
                app.Running.NavView = null;
            }

            // if we have a nav and highlight the nav item
            if ((app.Running.NavView) && (this.noNav.indexOf(Backbone.history.fragment) === -1)) {
                if (fragment === '')
                    fragment = config.DEFAULT_VIEW;

                app.Running.NavView.updateHighlight();
            }

            // render our page within the app
            app.Running.AppView.renderPage(app.Running.currentView);
        }
    });
})();
