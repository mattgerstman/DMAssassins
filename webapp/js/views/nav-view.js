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
		  this.listenTo(app.Running.TargetModel, 'change', this.handleTarget)
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
		if (!app.Running.TargetModel.get('user_id'))
		{
			  this.hideTarget();
			  return;
		}
		this.showTarget();
	  },
	  
	  // hides the target nav item
	  hideTarget: function(){
		  $('#nav_target').addClass('hide');
	  },
	  
	  // shows the target nav item
  	  showTarget: function(){
		  $('#nav_target').removeClass('hide');
	  }

  })
  
})(jQuery);