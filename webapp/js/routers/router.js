var app = app || {};

(function(){

	app.Router = Backbone.Router.extend({
		routes: {
			'' : 'index',
			'target' : 'target'
		},
		initialize: function() {
			app.session = app.session || {};
			app.session.email="Matt"
			this.appView = new app.AppView();
			this.navView = new app.NavView();
			this.appView.render();	
			this.navView.render();	
		},
		index : function() {
			console.log('index');
//			this.target()
		},
		target : function() {
			console.log('target');
			this.currentView = new app.UserView();
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