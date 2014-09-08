// handles the nav bar at the top
// js/views/nav-view.js

var app = app || {Models:{}, Views:{}, Routers:{}, Running:{}, Session:{}};

(function($){
 'use strict';
  app.Views.NavView = Backbone.View.extend({
	   
	    
	  template: _.template( $('#nav-template').html() ),
	  el: '#nav_body',
	  
	  tagName: 'nav',
	  
	  events: {
			'click li' : 'select'
	  },	  
	  
	  // constructor
	  initialize : function (){	  
	  	  this.model = app.Running.NavModel;
		  this.listenTo(this.model, 'change', this.render)
		  this.listenTo(app.Running.UserGamesModel, 'game-change', this.handleTarget)
	  },
	  
	  // if we don't have a target hide that view
	  render: function(){
		this.$el.html( this.template ( this.model.attributes ) );

		return this;  
	  },
	  
	  // select an item on the nav bar
	  select: function(event){
			var target = event.currentTarget;
			this.highlight(target)

	  },
	  
	  // highlight an item on the nav bar and unhighlight the rest of them
	  highlight: function(elem) {
	  	if ($(elem).hasClass('dropdown_parent')) {
		  	return;
	  	}
	  	
	  	if ($(elem).hasClass('dropdown_item')) {
	  		var dropdown = $(elem).attr('dropdown');
	  		var parent = '#'+dropdown+'_parent';	  		
	  		elem = parent;
	  	}
		$('.active').removeClass('active');
		$(elem).addClass('active');
	  },
	  
	  handleTarget: function(){
	  	var that = this;
	  
	  	app.Running.TargetModel.changeGame(app.Session.getGameId());		
		app.Running.TargetModel.fetch({
			success: function(model, error) {
				that.enableTarget();				
			},
			error: function(model, error) {
				if (error.status == 404){
					app.Running.TargetModel.set('user_id', null);					
				}
				that.disableTarget();
			}
		
		});
				
	  },
	  
	  // hides the target nav item
	  enableTarget: function(){
		  $('#nav_target').removeClass('disabled');
	  },
	  
	  // shows the target nav item
  	  disableTarget: function(){
		  $('#nav_target').addClass('disabled');
	  }

  })
  
})(jQuery);