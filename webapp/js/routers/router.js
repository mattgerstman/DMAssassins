// Handles all the URL *Magic*
// js/routers/router.js
var app = app || {Models:{}, Views:{}, Routers:{}, Running:{}, Session:{}};

(function(){

	app.Routers.Router = app.Routers.BaseRouter.extend({
	
		// Sometime we wanna go back and tis is the only wat to do that.
		history: [],
		// All the routes
		routes: {
			'' 				: 'target',
			'login' 		: 'login',
			'logout' 		: 'logout',
			'target' 		: 'target',
			'my_profile' 	: 'my_profile',
			'multigame' 	: 'multigame',
			'leaderboard' 	: 'leaderboard',
			'create_game'	: 'create_game',
			'join_game' 	: 'join_game',
			'rules' 		: 'rules',
			'switch_game' 	: 'switch_game'
			
		},
		// routes that require we have a game that has been started
		requiresGameStarted : ['#target', '#', ''],
		
		// routes that require just authentication
		requiresJustAuth : ['#multigame'],
		
		// routes that require we have a game and we're authenticated
		requiresGameAndAuth : ['#my_profile', '#join_game', '#leaderboard', '#rules'],
		
		// routes that should hide the nav bar
		noNav : ['login', 'multigame'],
			
		// routes that a logged in user can't access
		preventAccessWhenAuth : ['#login'],
		
		// place to redirect users for requiresGameAndAuth
		redirectWithoutGame : '#multigame',
		
		// place to redirect users who aren't logged in
		redirectWithoutAuth : '#login',
		
		before : function(params, next){

			// is the user authenticated
			var isAuth = app.Session.get('authenticated');
			var path = Backbone.history.location.hash;
			
			// do we need a game and authentication
			var needGameAndAuth = _.contains(this.requiresGameAndAuth, path);
			
			// do we need authentication
			var needAuth 		= _.contains(this.requiresJustAuth, path);

			// should we prevent this page if authorized
			var cancelAccess 	= _.contains(this.preventAccessWhenAuth, path);
			
			// does this route need a running game
			var needStarted 	= _.contains(this.requiresGameStarted, path);
			
			// is there a game
			var isGame   		= app.Session.get('game') !== null;
			
			// is the game started
			var isStarted 		= app.Session.get('game') && (app.Session.get('game').game_started);

/*
			Variables I use when shit's not routing properly /**//*/
			console.log('path:', path);
			console.log('needGameAndAuth: ', needGameAndAuth);
			console.log('isGame: ', isGame);
			console.log('isStarted: ', isStarted);
			console.log('needAuth: ', needAuth);
			console.log('isAuth: ', isAuth);
			console.log('cancelAccess: ', cancelAccess);
/**/

			// Do we need authentication
			if((needAuth || needGameAndAuth) && !isAuth){
				app.Session.set('redirect_from', path);
				Backbone.history.navigate(this.redirectWithoutAuth, { trigger : true });
			}
			// do we need authentication and a game
			else if (needGameAndAuth && !isGame) {
				Backbone.history.navigate(this.redirectWithoutGame, { trigger : true });
			}
			// do we need a game and is it started
			else if (needStarted && !isStarted) {
				return;
			}
			// are they logged in and trying to hit the login page
			else if(isAuth && cancelAccess) {
				Backbone.history.navigate('', { trigger : true });			
			// nothing is wrong! let them pass.	
			} else {
				//No problem handle the route
				return next();
			}			
		},
		// called after we're done routing, unused but build into the baserouter so we're leaving it
		after: function(){
			this.history.push(Backbone.history.fragment);
		},
		// go to the previous
		back: function(){
			this.history.pop();
			history.back();
		},
		// login route
		login : function() {
			app.Running.currentView = new app.Views.LoginView();
			this.render();				
		},
		// logout route
		logout : function() {
			app.Session.clear()
			this.navigate('login', true)
		},
		// game selection route
		multigame : function() {
			app.Running.currentView = new app.Views.SelectGameView();
			app.Running.currentView.model.fetch();
			this.render();
		},
		// target route
		target : function() {
			//console.log('target');
			app.Running.currentView = new app.Views.TargetView();
			app.Running.TargetModel.changeGame(app.Session.getGameId());
			app.Running.currentView.model.fetch();
			this.render();
		},
		// create a new game route
		create_game: function() {
			app.Running.currentView = new app.Views.SelectGameView();
			this.render();
			app.Running.currentView.showCreateGame();		
		},
		// join a new game route
		join_game: function() {
			app.Running.currentView = new app.Views.SelectGameView();
			this.render();
			app.Running.currentView.loadJoinGame(app.Session.get('user_id'));
		},
		// profile route
		my_profile : function() {
			//console.log('profile');			
			app.Running.currentView = new app.Views.ProfileView();
			app.Running.currentView.model.changeGame();
			app.Running.currentView.model.fetch();
			this.render();
		},
		// leaderboard route
		leaderboard: function(){
			//console.log('leaderboard');			
			app.Running.currentView = new app.Views.LeaderboardView();
			app.Running.currentView.model.loadGame();
			app.Running.currentView.model.fetch();
			this.render();
		},
		// rules route
		rules : function() {
			//console.log('rules');			
			app.Running.currentView = new app.Views.RulesView();
			app.Running.currentView.model.loadGame();
			app.Running.currentView.model.fetch();
			this.render();
		},
		preventSwitchGameBack : ['join_game', 'create_game'],
		switch_game: function() {
			var lastFragment = this.history[this.history.length - 1];
			if (lastFragment === undefined || _.contains(this.preventSwitchGameBack, lastFragment)) {
				Backbone.history.navigate('my_profile', { trigger : true });	
				return;
			}
			
			this.back();
			
		},
		// render function, also determines weather or not to render the nav
		render : function(){
			var fragment = Backbone.history.fragment;
			// if it's a view with a nav and we don't have one, make one
			if ((this.noNav.indexOf(Backbone.history.fragment) == -1) && (fragment != 'login') && (!app.Running.navView))
			{
				app.Running.navView = new app.Views.NavView();
				app.Running.navView = app.Running.navView.render();
				app.Running.NavGameView = new app.Views.NavGameView(app.Session.get('user_id'));
				app.Running.NavGameView.model.fetch();
				app.Running.NavGameView.render();
			}
			// if it explicitely shouldn't have a nav and we have one kill it
			else if ((this.noNav.indexOf(Backbone.history.fragment) != -1) && (app.Running.navView))
			{
				app.Running.navView.$el.html('');
				app.Running.navView = null;
			}
			// if we have a nav and highlight the nav item
			if ((app.Running.navView) && (this.noNav.indexOf(Backbone.history.fragment) == -1))
			{
				//console.log(fragment);
				if (fragment === '')
					fragment = 'target';
					
				app.Running.navView.highlight('#nav_'+fragment)
				app.Running.NavGameView.updateText();
			}
			
			// render our page within the app
			app.Running.appView.renderPage(app.Running.currentView)			
		}		
	})

})()