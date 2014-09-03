// handles the game dropdown in the nav
// js/views/nav-game-view.js

var app = app || {Models:{}, Views:{}, Routers:{}, Running:{}, Session:{}};

(function($){
 'use strict';
  app.Views.NavGameView = Backbone.View.extend({
	   
	     
	  template: _.template( $('#nav-game-template').html() ),
	  el: '#games_dropdown',
	  
	  tagName: 'ul',
	  
	  events: {
			'click li a.switch_game' : 'select'
	  },	  
	  // constructor, loads a user id so we can get their games from the model
	  initialize : function (user_id) {	  
		  this.model = app.Running.GamesModel;
		  app.Running.GamesModel.loadUser(user_id);
		  this.listenTo(this.model, 'change', this.render)
		  this.listenTo(this.model, 'game-change', this.render);

	  },
	  
	  // loads the items into the dropdown and changes the dropdown title to the current game
	  render: function(){
		  this.$el.html( this.template ( this.model.attributes ) );
		  $('#games_header').text(app.Session.get('game').game_name);	
		  var game_id = app.Session.get('game').game_id;
		  $('#nav_'+game_id).addClass('hide');
	  },
	  	  
	  showJoinGame: function() {
		  var game_id = app.Session.get('game').game_id;
		  $('#nav_'+game_id).removeClass('hide');
		  $('#games_header').text('Join New Game');
		  
	  },
	  
	  	  
	  showCreateGame: function() {
		  var game_id = app.Session.get('game').game_id;
		  $('#nav_'+game_id).removeClass('hide');
		  $('#games_header').text('Create New Game');
		  
	  },	  
	  // select a game from the dropdown
	  select: function(event){
	  	var game_id = $(event.target).attr('game_id');	
	  	this.model.switchGame(game_id)

	  }
  })
  
})(jQuery);