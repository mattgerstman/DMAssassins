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
			if (localStorage.getItem('logged_in') === "true")
			{
				this.navigate('target');
			}
			else
			{	
				this.navigate('index');
			}		
		},
		index : function() {

			app.Running.currentView = new app.Views.LoginView();
			app.Running.currentView.render()				
				
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
			app.Running.currentView.model.fetch();
		}/*
,
		reload : function(){
			console.log('reload');
			if (Backbone.history.fragment != '') {
				app.Running.currentView.model.changeUser(app.Session.username);	
			}
			
		}
*/
	})

})()