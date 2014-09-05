// handles game selection
// js/views/select-game-view.js

var app = app || {Models:{}, Views:{}, Routers:{}, Running:{}, Session:{}};

(function($){
 'use strict';
  app.Views.SelectGameView = Backbone.View.extend({
	   
	     
	  template: _.template( $('#select-game-template').html() ),
	  tagName: 'div',
	  events: {
			'click .show-create-game' 			: 'showCreateGame',
			'click .show-join-game'   			: 'showJoinGame',
			'click .create-game-submit'			: 'createGame',
			'click .join-game-submit'			: 'joinGame',
			'click .create-or-join-back'		: 'goBack',
			'click #create_game_need_password'	: 'togglePassword',
			'change #games'						: 'checkPassword'
			
	  },
	  // previous page, may depricate
	  loaded_from : 'login',
	  // constructor
	  initialize : function (params){
		  this.model = app.Running.SelectGamesModel;
		  this.listenTo(this.model, 'change', this.render)
		  this.listenTo(this.model, 'fetch', this.render)	
		  this.listenTo(this.model, 'finish_set_game', this.finish)	  
	  },
	  // shows the create game subview
	  showCreateGame: function(){
	  		$('.logo').addClass('hide');
			$('#create-or-join').addClass('hide');
			$('#create-game').addClass('select-game-active');
			$('#create-game').removeClass('hide');
	  },
	  // shows the join game subview
	  showJoinGame: function(){
			$('.logo').addClass('hide');
			$('#create-or-join').addClass('hide');
			$('#join-game').addClass('select-game-active');
			$('#join-game').removeClass('hide');
	  },
	  // cancels the game creation/selection
	  goBack: function(){
		  if (app.Session.get('authenticated'))
		  {
		  		// im sure theres a back function, find it
		  		history.back();
		  		return;
		  }
		  $('.select-game-active').addClass('hide').removeClass('select-game-active');
		  $('#create-or-join').removeClass('hide');
	  },
	  // show the create game s ubview
	  createGame: function(){
	  		var name = $('#create_game_name').val();
	  		var password = null;
	  		if ($('#create_game_need_password').is(':checked'))
	  		{
		  		password = $('#create_game_password').val();
		  		
	  		}
			this.model.create(name, password);
	  },
	  // loads the join game later view
	  loadJoinGame: function(user_id){
	  	var that = this;
	  	that.showJoinGame();
	  	this.model.setUser(user_id);
	  	this.model.fetch({
	  		success: function(){
				that.showJoinGame();  	
	  		}
	  	});

	  },
	  // posts to the join game model
	  joinGame: function(){
	  		console.log('join');
	  		var user_id = app.Session.get('user_id');
		  	var game_id = $('#games option:selected').val();
	  		var password = $('#join_game_password').val();
			this.model.join(game_id, user_id, password);
	  },
	  // finish up and navigate to your profile
	  finish: function(){
	  		app.Running.UserGamesModel.switchGame(app.Session.get('game').game_id);
			Backbone.history.navigate('my_profile', { trigger : true });  
	  },
	  // toggles the password entry field on create game
	  togglePassword: function(e){			
			$('#create_game_password').attr('disabled', !e.target.checked);
	  },
	  // toggles the password entry field on join game
	  checkPassword: function(e){
	  	  	var need_password = $('#games option:selected').attr('game_has_password') == 'true';
	  	  	$('#join_game_password').attr('disabled', !need_password);	  	  	
		  
	  },
	  render: function(){
		this.$el.html( this.template ( this.model.attributes ) );
		return this;  
	  }	    
 
  })
  
})(jQuery);