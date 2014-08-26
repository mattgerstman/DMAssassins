  // js/views/select-game-view.js

var app = app || {Models:{}, Views:{}, Routers:{}, Running:{}, Session:{}};

(function($){
 'use strict';
  app.Views.SelectGameView = Backbone.View.extend({
	   
	     
	  template: _.template( $('#select-game-template').html() ),
	  tagName: 'div',
	  events: {
			'click .create-game' : 'create_game',
			'click .join-game'   : 'join_game'			
	  },
	  initialize : function (params){
	  	this.model = app.Running.GamesModel;

		  this.listenTo(this.model, 'change', this.render)
		  this.listenTo(this.model, 'fetch', this.render)		  
	  },
	  create_game: function(){
			$('#create-or-join').addClass('hide');
			$('#create-game').removeClass('hide');
	  },
	  join_game: function(){
			$('#create-or-join').addClass('hide');
			$('#join-game').removeClass('hide');
	  },
	  render: function(){
//	  	this.$el.hide()
		this.$el.html( this.template ( this.model.attributes ) );
		this.model.fetch();
//		this.$el.fadeIn(250);
		return this;  
	  }	    
 
  })
  
})(jQuery);