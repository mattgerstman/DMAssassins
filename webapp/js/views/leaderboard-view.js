  // js/views/profile-view.js

var app = app || {Models:{}, Views:{}, Routers:{}, Running:{}, Session:{}};

(function($){
 'use strict';
  app.Views.LeaderboardView = Backbone.View.extend({
	   
	     
	  template: _.template( $('#leaderboard-template').html() ),
	  tagName: 'div',
	  
	  initialize : function (params){
	  	this.model = app.Running.LeaderboardModel;

		  this.listenTo(this.model, 'change', this.render)
		  this.listenTo(this.model, 'fetch', this.render)		  
	  },
	  
	  render: function(){
//	  	this.$el.hide()
		this.$el.html( this.template ( this.model.attributes ) );
//		this.$el.fadeIn(250);
		return this;  
	  }	    
 
  })
  
})(jQuery);