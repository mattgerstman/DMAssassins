// js/views/nav-game-view.js

var app = app || {Models:{}, Views:{}, Routers:{}, Running:{}, Session:{}};

(function($){
 'use strict';
  app.Views.NavGameView = Backbone.View.extend({
	   
	     
	  template: _.template( $('#nav-game-template').html() ),
	  el: '#game_dropdown',
	  
	  tagName: 'ul',
	  
	  events: {
			'click li a.switch_game' : 'select'
	  },	  
	  initialize : function (user_id){	  
	  	  this.model = app.Running.GamesModel;
	  	  app.Running.GamesModel.loadUser(user_id);
		  this.listenTo(this.model, 'change', this.render)
		  this.listenTo(this.model, 'game-change', this.render);

	  },
	  
	  render: function(){
			this.$el.html( this.template ( this.model.attributes ) );			

			$('#game_header').text(app.Session.get('game').game_name);	
			var game_id = app.Session.get('game_id');
			$('#nav_'+game_id).addClass('hide');
			

	  },
	  select: function(event){
	  	var game_id = $(event.target).attr('game_id');
	  	this.model.switchGame(game_id)
//			var target = event.currentTarget;
//			this.highlight(target)

	  }
  })
  
})(jQuery);