  // js/views/profile-view.js

var app = app || {Models:{}, Views:{}, Routers:{}, Running:{}, Session:{}};

(function($){
 'use strict';
  app.Views.ProfileView = Backbone.View.extend({
	   
	     
	  template: _.template( $('#profile-template').html() ),
	  tagName: 'div',
	  
        // The DOM events specific to an item.
		events: {
	      'click .thumbnail': 'showFullImage'
	    },
	  
	  showFullImage: function(){
		  $('#photoModal').modal()  
	  },
	  
	  initialize : function (params){
	  	this.model = app.Running.ProfileModel;

		  this.listenTo(this.model, 'change', this.render)
		  this.listenTo(this.model, 'fetch', this.render)		  
	  },
	  
	  render: function(){
//	  	this.$el.hide()
		this.$el.html( this.template ( this.model.attributes ) );
//		this.$el.fadeIn(250);
		return this;  
	  }	    
 
  })
  
})(jQuery);