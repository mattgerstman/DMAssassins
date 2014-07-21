  // js/views/user.js

var app = app || {};

(function($){
 'use strict';
  app.UserView = Backbone.View.extend({
	   
	     
	  template: _.template( $('#user-template').html() ),
	  
	  tagName: 'div',
	  
        // The DOM events specific to an item.
		events: {
	      'click .thumbnail': 'showFullImage'
	    },
	  
	  showFullImage: function(){
		  $('#photoModal').modal()  
	  },
	  
	  model: new app.User(),
	  
	  initialize : function (){
	  
		  this.listenTo(this.model, 'change', this.render)
	  },
	  
	  render: function(){
		this.$el.html( this.template ( this.model.attributes ) );
		return this;  
	  }	    
  })
  
})(jQuery);