var app = app || {Models:{}, Views:{}, Routers:{}, Running:{}, Session:{}};

(function(){

	app.Routers.Router = Backbone.Router.extend({
		routes: {
			'' : 'index',
			'target' : 'target',
			'my_profile' : 'my_profile'
		},
		initialize: function() {
			app.Running.appView = new app.Views.AppView();
			app.Running.appView.render();				
			this.navigate('index')
		},
		index : function() {
			if (app.Session.user_id === undefined)
			{
				app.Running.currentView = new app.Views.LoginView();
				app.Running.currentView.render()
			}
			else
			{
				this.target();
			}
				
		},
		target : function() {
			console.log('target');
			app.Running.currentView = new app.Views.TargetView();
			this.render();
			app.Running.navView.highlight('#nav_target')
		},
		my_profile : function() {
			console.log('profile');			
			app.Running.currentView = new app.Views.ProfileView();
			this.render();
			app.Running.navView.highlight('#nav_my_profile')
		},
		render : function(){
			console.log('fragment');
			console.log(Backbone.history.fragment);
			if ((Backbone.history.fragment != '') && (app.Running.navView === undefined))
			{
				app.Running.navView = new app.Views.NavView();
				app.Running.navView = app.Running.navView.render();
			}
			app.Running.appView.renderPage(app.Running.currentView)
		}
		
	})

})()