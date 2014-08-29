  // js/views/select-game-view.js

var app = app || {Models:{}, Views:{}, Routers:{}, Running:{}, Session:{}};

(function($){
 'use strict';
  app.Views.SelectGameView = Backbone.View.extend({
	   
	     
	  template: _.template( $('#select-game-template').html() ),
	  tagName: 'div',
	  events: {
			'click .show-create-game' 			: 'show_create_game',
			'click .show-join-game'   			: 'show_join_game',
			'click .create-game-submit'			: 'create_game',
			'click .join-game-submit'			: 'join_game',
			'click .create-or-join-back'		: 'go_back',
			'click #create_game_need_password'	: 'toggle_password',
			'change #games'						: 'check_password'
			
	  },
	  loaded_from : 'login',
	  initialize : function (params){
	  	this.model = app.Running.GamesModel;

		  this.listenTo(this.model, 'change', this.render)
		  this.listenTo(this.model, 'fetch', this.render)	
		  this.listenTo(this.model, 'finish_set_game', this.finish)	  
	  },
	  show_create_game: function(){
			$('#create-or-join').addClass('hide');
			$('#create-game').addClass('select-game-active');
			$('#create-game').removeClass('hide');
	  },
	  show_join_game: function(){
			$('#create-or-join').addClass('hide');
			$('#join-game').addClass('select-game-active');
			$('#join-game').removeClass('hide');
	  },
	  go_back: function(){
		  if (this.loaded_from != 'login')
		  {
			  Backbone.history.navigate(this.loaded_from, { trigger : true });
			  return;
		  }
		  $('.select-game-active').addClass('hide').removeClass('select-game-active');
		  $('#create-or-join').removeClass('hide');
	  },
	  create_game: function(){
	  		var name = $('#create_game_name').val();
	  		var password = null;
	  		if ($('#create_game_need_password').is(':checked'))
	  		{
		  		password = $('#create_game_password').val();
		  		
	  		}
			this.model.create(name, password);
	  },
	  load_join_game: function(loaded_from){
	  	this.loaded_from = loaded_from;
	  	this.model.set('user_id', null);
		this.show_join_game();	

	  },
	  join_game: function(){
	  		console.log('join');
	  		var user_id = app.Session.get('user_id');
		  	var game_id = $('#games option:selected').val();
	  		var password = $('#join_game_password').val();
			this.model.join(game_id, user_id, password);
	  },
	  finish: function(){
			Backbone.history.navigate('my_profile', { trigger : true });  
	  },
	  toggle_password: function(e){			
			$('#create_game_password').attr('disabled', !e.target.checked);
	  },
	  check_password: function(e){
	  	  	var need_password = $('#games option:selected').attr('game_has_password') == 'true';
	  	  	$('#join_game_password').attr('disabled', !need_password);	  	  	
		  
	  },
	  render: function(){
//	  	this.$el.hide()
		this.$el.html( this.template ( this.model.attributes ) );
		//this.model.fetch();
//		this.$el.fadeIn(250);
		return this;  
	  }	    
 
  })
  
})(jQuery);