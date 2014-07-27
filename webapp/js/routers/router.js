var app = app || {Models:{}, Views:{}, Routers:{}, Running:{}};

(function(){

	app.Routers.Router = Backbone.Router.extend({
		routes: {
			'' : 'index',
			'target' : 'target'
		},
		initialize: function() {
			app.session = app.session || {};
			app.session.email="Matt"
			this.appView = new app.Views.AppView();
			this.appView.render();	
//			this.navView = new app.NavView();
//			this.navView.render()

		},
		index : function() {
			this.currentView = new app.Views.LoginView();
			this.currentView.render();
//			this.target()
		},
		target : function() {
			console.log('target');
			this.navView = new app.Views.NavView();
			this.navView.render();
			this.currentView = new app.Views.TargetView({'username' : 'Matt'});
			this.currentView.render()
			console.log(this.currentView.model);
			this.currentView.model.fetch({
				success: function (data) {
	                console.log(data);
	                // Note that we could also 'recycle' the same instance of EmployeeFullView
	                // instead of creating new instances
					}
				}
            )
			this.render();
			this.navView.highlight('#nav_target')
		},
		render : function(){
			this.appView.renderPage(this.currentView)
		}
		
	})

})()