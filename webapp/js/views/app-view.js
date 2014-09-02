// loads pages within the body of the app
// js/views/app-view
var app = app || {Models:{}, Views:{}, Routers:{}, Running:{}, Session:{}};

(function($){
	'use strict';
	app.Views.AppView = Backbone.View.extend({
			
		el: '#app',
		// constructor
		initialize: function(){
			this.$body = $('#app_body');

		},
		// renders a page within the body of the app
		renderPage : function(page) {
			this.$body.html(page.render().el);
		}
	})
})(jQuery);