  // js/views/select-game-view.js

var app = app || {Models:{}, Views:{}, Routers:{}, Running:{}, Session:{}};

(function($){
 'use strict';
  app.Views.SelectGameView = Backbone.View.extend({
	   
	     
	  template: _.template( $('#select-game-template').html() ),
	  tagName: 'div',
	  
	  initialize : function (params){
	  	this.model = app.Running.GamesModel;

		  this.listenTo(this.model, 'change', this.render)
		  this.listenTo(this.model, 'fetch', this.render)		  
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