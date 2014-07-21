var app = app || {};

(function($){
	'use strict';
	app.AppView = Backbone.View.extend({
		
		
		el: '#app',
		events : {
			'click #nav_target' : 'target'	
		},

		initialize: function(){
			this.$body = $('#app_body');

		},
		target : function(){
			var userView = new app.UserView();
			this.$body.html(userView.render().el);			
		}
	})
})(jQuery);