var app = app || {};

(function($){
	'use strict';
	app.AppView = Backbone.View.extend({
		
		
		el: '#app',

		initialize: function(){
			this.$body = $('#app_body');

		},
		renderPage : function(page) {
			this.$body.html(page.render().el);
		}
	})
})(jQuery);