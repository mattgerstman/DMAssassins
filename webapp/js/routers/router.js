var app = app || {};

(function(){

	app.Router = Backbone.Router.extend({
		routes: {
			'' : 'index',
			'target' : 'target'
		},
		index : function() {
			console.log('index');
			this.target()
		},
		target : function() {
			console.log('target');
			Backbone.trigger('click #nav_target');
		}
		
	})

})()