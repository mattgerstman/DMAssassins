var app = app || {Models:{}, Views:{}, Routers:{}, Running:{}};

(function(){

	app.Routers.Router = Backbone.Router.extend({
		routes: {
			'' : 'index',
			'target' : 'target',
			'my_profile' : 'my_profile'
		},
		initialize: function() {
			this.appView = new app.Views.AppView();
			this.appView.render();	
		},
		index : function() {
			this.currentView = new app.Views.LoginView();
			this.currentView.render();			
			//this.target()
		},
		target : function() {
			this.navView = new app.Views.NavView();
			this.navView.render()
			console.log('target');
			this.currentView = new app.Views.TargetView({'username' : 'Jimmy'});
			this.render();
			this.navView.highlight('#nav_target')
		},
		my_profile : function() {
			console.log('profile');			
			this.currentView = new app.Views.ProfileView({'username' : 'MattGerstman'});
			this.render();
			this.navView.highlight('#nav_my_profile')
		},
		render : function(){
			this.appView.renderPage(this.currentView)
		}
		
	})

})()