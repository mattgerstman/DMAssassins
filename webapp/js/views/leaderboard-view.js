// shows the list of high scores
// js/views/leaderboard-view.js

var app = app || {Models:{}, Views:{}, Routers:{}, Running:{}, Session:{}};

(function($){
 'use strict';
  app.Views.LeaderboardView = Backbone.View.extend({
	   
	     
	  template: _.template( $('#leaderboard-template').html() ),
	  tagName: 'div',
	  
	  // constructor
	  initialize : function (params){
	  	this.model = app.Running.LeaderboardModel;

		  this.listenTo(this.model, 'change', this.render)
		  this.listenTo(this.model, 'fetch', this.render)		  
	  },
	  // renderer
	  render: function(){
		this.$el.html( this.template ( this.model.attributes ) );
		this.$el.find('#leaderboard_table').dataTable();		
		return this;  
	  }	    
 
  })
  
})(jQuery);