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
		var numCols = this.model.get('teams_enabled') + 2;
		var options = {
			 paging:		false,
			 searching: 	false,
			 info: 		    false,
			 order:			[[numCols-1, 'desc'], [numCols, 'desc']]
		};

		this.$el.find('#user_leaderboard_table').dataTable(options);
		
		if (this.model.get('teams_enabled'))
		{
			options.order = [[4, 'desc']]
			this.$el.find('#team_leaderboard_table').dataTable(options);
		}


		

		return this;  
	  }	    
 
  })
  
})(jQuery);