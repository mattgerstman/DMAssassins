var app = app || {Models:{}, Views:{}, Routers:{}, Running:{}, Session:{}};

(function(){

	app.Routers.Router = app.Routers.BaseRouter.extend({
		routes: {
			'' : 'target',
			'login' : 'login',
			'logout' : 'logout',
			'target' : 'target',
			'my_profile' : 'my_profile',
			'multigame' : 'multigame',
			'leaderboard' : 'leaderboard',
			'join_game' : 'join_game',
			'rules' : 'rules'
			
		},
		
		// Routes that need authentication and if user is not authenticated
		// gets redirect to login page
		requiresGameStarted : ['#target', '#', ''],
		requiresJustAuth : ['#multigame'],
		requiresGameAndAuth : ['#my_profile', '#join_game'],
		noNav : ['login', 'multigame'],
		// Routes that should not be accessible if user is authenticated
		// for example, login, register, forgetpasword ...
		preventAccessWhenAuth : ['#login'],
		redirectWithoutGame : '#multigame',
		redirectWithoutAuth : '#login',
		redirectWithoutGameStarted : 'my_profile',
		before : function(params, next){
			//Checking if user is authenticated or not
			//then check the path if the path requires authentication 
			var isAuth = app.Session.get('authenticated');
			var path = Backbone.history.location.hash;
			//console.log('path');
			//console.log(path);
			var needGameAndAuth = _.contains(this.requiresGameAndAuth, path);
			var needAuth 		= _.contains(this.requiresJustAuth, path);

			var cancelAccess 	= _.contains(this.preventAccessWhenAuth, path);
			var needStarted 	= _.contains(this.requiresGameStarted, path);
			var isGame   		= app.Session.get('game_id') !== null;
			var isStarted 		= app.Session.get('game') && (app.Session.get('game').game_started === 'true');

/*
			console.log('path:', path);
			console.log('needGameAndAuth: ', needGameAndAuth);
			console.log('isGame: ', isGame);
			console.log('isStarted: ', isStarted);
			console.log('needAuth: ', needAuth);
			console.log('isAuth: ', isAuth);
			console.log('cancelAccess: ', cancelAccess);
*/



			if((needAuth || needGameAndAuth) && !isAuth){
				//If user gets redirect to login because wanted to access
				// to a route that requires login, save the path in session
				// to redirect the user back to path after successful login
				app.Session.set('redirectFrom', path);
				Backbone.history.navigate(this.redirectWithoutAuth, { trigger : true });
			}
			else if (needGameAndAuth && !isGame) {
				Backbone.history.navigate(this.redirectWithoutGame, { trigger : true });
			}
			else if (needStarted && !isStarted) {
				Backbone.history.navigate(this.redirectWithoutGameStarted, { trigger : true });
			}
			else if(isAuth && cancelAccess){
				//User is authenticated and tries to go to login, register ...
				// so redirect the user to home page
				Backbone.history.navigate('', { trigger : true });
			}else{
				//No problem handle the route
				return next();
			}			
		},

		after : function(){
			//empty
		},
		initialize: function() {
			if (app.Session.get('authenticated') !== true)
			{
				//console.log('not authenticated');
				//this.navigate('login', true)
			}
		},	
		login : function() {
			
			app.Running.currentView = new app.Views.LoginView();
			this.render();
				
		},
		logout : function() {
			app.Session.clear()
			this.navigate('login', true)
		},
		multigame : function() {
			app.Running.currentView = new app.Views.SelectGameView();
			app.Running.currentView.model.fetch();
			this.render();
		},
		target : function() {
			//console.log('target');
			app.Running.currentView = new app.Views.TargetView();
			app.Running.TargetModel.changeUser(app.Session.get('user_id'))
			app.Running.currentView.model.fetch();
			this.render();
		},
		join_game: function() {
			app.Running.currentView = new app.Views.SelectGameView();
			this.render();
			app.Running.currentView.load_join_game('target');				
		},
		my_profile : function() {
			//console.log('profile');			
			app.Running.currentView = new app.Views.ProfileView();
			app.Running.currentView.model.fetch();
			this.render();
		},
		leaderboard: function(){
			//console.log('leaderboard');			
			app.Running.currentView = new app.Views.LeaderboardView();
			app.Running.currentView.model.fetch();
			this.render();
		},
		rules : function() {
			//console.log('rules');			
			app.Running.currentView = new app.Views.RulesView();
			app.Running.currentView.model.fetch();
			this.render();
		},
		render : function(){
			var fragment = Backbone.history.fragment;
			if ((this.noNav.indexOf(Backbone.history.fragment) == -1) && (fragment != 'login') && (!app.Running.navView))
			{
				app.Running.navView = new app.Views.NavView();
				app.Running.navView = app.Running.navView.render();
				app.Running.NavGameView = new app.Views.NavGameView(app.Session.get('user_id'));
				app.Running.NavGameView.model.fetch();
				app.Running.NavGameView.render();
			}
			else if ((this.noNav.indexOf(Backbone.history.fragment) != -1) && (app.Running.navView))
			{
				app.Running.navView.$el.html('');
				app.Running.navView = null;
			}
			if ((app.Running.navView) && (this.noNav.indexOf(Backbone.history.fragment) == -1))
			{
				//console.log(fragment);
				if (fragment === '')
					fragment = 'target';
					
				app.Running.navView.highlight('#nav_'+fragment)	
			}
			app.Running.appView.renderPage(app.Running.currentView)			
		},
		fetchError : function(error){
			//If during fetching data from server, session expired
			// and server send 401, call getAuth to get the new CSRF
			// and reset the session settings and then redirect the user
			// to login
			if(error.status === 401){
				Session.getAuth(function(){
					Backbone.history.navigate('login', { trigger : true });
				});
			}
		}
		
	})

})()