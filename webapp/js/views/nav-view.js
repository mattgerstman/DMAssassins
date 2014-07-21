// js/views/nav-view.js

var app = app || {};

(function($){
 'use strict';
  app.NavView = Backbone.View.extend({
	   
	     
	  template: _.template( $('#nav-template').html() ),
	  el: '#nav_body',
	  
	  tagName: 'nav',
	  
	  events: {
			'click li' : 'select'
	  },	  
	  model: new app.Nav(),
	  
	  initialize : function (){	  
		  this.listenTo(this.model, 'change', this.render)
	  },
	  
	  render: function(){
		this.$el.html( this.template ( this.model.attributes ) );
		return this;  
	  },
	  select: function(event){
			var target = event.currentTarget;
			console.log($(target));
			$('.active').removeClass('active');
			$(target).addClass('active');
			
	  }
  })
  
})(jQuery);