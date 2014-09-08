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
		  this.model = app.Running.UserGamesModel;
		  this.model.loadUser(user_id);
		  this.listenTo(this.model, 'change', this.render)
		  this.listenTo(this.model, 'game-change', this.render);

	  },
	  showCurrentGame: function() {
  		  var game_id = app.Session.getGameId();
		  $('#nav_'+game_id).removeClass('hide');
	  },
	  updateText:function() {
		  if (Backbone.history.fragment == 'join_game')
		  {
		  	this.showCurrentGame();
			$('#games_header').text('Join Game');
			$('#games_header_short').text('Join Game');
			return;
		  }
		  
		  if (Backbone.history.fragment == 'create_game')
		  {
		  	this.showCurrentGame();
			$('#games_header').text('Create Game');
			$('#games_header_short').text('Create Game');
			return;
		  }
		  
		  var game_name = app.Session.get('game').game_name;
		  $('#games_header').text(game_name);
	  	  var max = 9;	  
		  if (game_name.length > max)
		  {
				game_name = game_name.substr(0,max-3) + '...';	  
		  }		  	
		  $('#games_header_short').text(game_name);
		  
  		  var game_id = app.Session.getGameId();
		  $('#nav_'+game_id).addClass('hide');
	  },
	  // loads the items into the dropdown and changes the dropdown title to the current game
	  render: function(){
		  this.$el.html( this.template ( this.model.attributes ) );
		  this.updateText();
	  },	  
	  // select a game from the dropdown
	  select: function(event){
	  	var game_id = $(event.target).attr('game_id');
	  	this.model.switchGame(game_id)

	  }
  })
  
})(jQuery);