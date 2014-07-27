  // js/views/user-view.js

var app = app || {Models:{}, Views:{}, Routers:{}, Running:{}};

(function($){
 'use strict';
  app.Views.TargetView = Backbone.View.extend({
	   
	     
	  template: _.template( $('#user-template').html() ),
	  
	  tagName: 'div',
	  
        // The DOM events specific to an item.
		events: {
	      'click .thumbnail': 'showFullImage'
	    },
	  
	  showFullImage: function(){
		  $('#photoModal').modal()  
	  },
	  
	  initialize : function (params){
	  	this.model = new app.Models.User(params)
		this.model.fetch()
		  this.listenTo(this.model, 'change', this.render)
		  this.listenTo(this.model, 'fetch', this.render)		  
	  },
	  
	  render: function(){
		this.$el.html( this.template ( this.model.attributes ) );
		return this;  
	  }	    
  })
  
})(jQuery);