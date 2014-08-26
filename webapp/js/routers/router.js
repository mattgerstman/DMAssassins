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
			'rules' : 'rules'
			
		},
		
		// Routes that need authentication and if user is not authenticated
		// gets redirect to login page
		requiresJustAuth : ['#multigame'],
		requiresGameAndAuth : ['#my_profile', '#target', '#', ''],
		noNav : ['login', 'multigame'],
		// Routes that should not be accessible if user is authenticated
		// for example, login, register, forgetpasword ...
		preventAccessWhenAuth : ['#login'],
		redirectWithoutGame : '#multigame',
		redirectWithoutAuth : '#login',
		before : function(params, next){
			//Checking if user is authenticated or not
			//then check the path if the path requires authentication 
			var isAuth = app.Session.get('authenticated');
			var path = Backbone.history.location.hash;
			console.log('path');
			console.log(path);
			var needGame = _.contains(this.requiresGameAndAuth, path);
			var isGame   = app.Session.get('game_id');
			var needAuth = _.contains(this.requiresAuthJustAuth, path);
			var cancelAccess = _.contains(this.preventAccessWhenAuth, path);
			
			console.log('needGame: ', needGame);
			console.log('isGame: ', isGame);
			
			if (needGame && (isGame == 'null')) {
				Backbone.history.navigate(this.redirectWithoutGame, { trigger : true });
			}
			else if(needAuth && !isAuth){
				//If user gets redirect to login because wanted to access
				// to a route that requires login, save the path in session
				// to redirect the user back to path after successful login
				app.Session.set('redirectFrom', path);
				Backbone.history.navigate(this.redirectWithoutAuth, { trigger : true });
			}else if(isAuth && cancelAccess){
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
				console.log('not authenticated');
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
			this.render();
		},
		target : function() {
			console.log('target');
			app.Running.currentView = new app.Views.TargetView();
			app.Running.TargetModel.changeUser(app.Session.get('user_id'))
			app.Running.currentView.model.fetch();
			this.render();
		},
		my_profile : function() {
			console.log('profile');			
			app.Running.currentView = new app.Views.ProfileView();
			app.Running.currentView.model.fetch();
			this.render();
		},
		leaderboard: function(){
			console.log('leaderboard');			
			app.Running.currentView = new app.Views.LeaderboardView();
			app.Running.currentView.model.fetch();
			this.render();
		},
		rules : function() {
			console.log('rules');			
			app.Running.currentView = new app.Views.RulesView();
			app.Running.currentView.model.fetch();
			this.render();
		},
		render : function(){
			var fragment = Backbone.history.fragment;
			if ((this.noNav.indexOf(Backbone.history.fragment) == -1) && (fragment != 'login') && (app.Running.navView === undefined))
			{
				app.Running.navView = new app.Views.NavView();
				app.Running.navView = app.Running.navView.render();
			}
			else if ((this.noNav.indexOf(Backbone.history.fragment) != -1) && (app.Running.navView !== undefined))
			{
				app.Running.navView.remove();
			}
			if ((app.Running.navView !== undefined) && (this.noNav.indexOf(Backbone.history.fragment) == -1))
			{
				console.log(fragment);
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