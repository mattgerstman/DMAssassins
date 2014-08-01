  // js/views/user-view.js

var app = app || {Models:{}, Views:{}, Routers:{}, Running:{}, Session:{}};

(function($){
 'use strict';
  app.Views.TargetView = Backbone.View.extend({
	   
	     
	  template: _.template( $('#target-template').html() ),
	  
	  tagName: 'div',
	  
        // The DOM events specific to an item.
		events: {
	      'click .thumbnail': 'showFullImage',
	      'click #kill' : 'kill'
	    },
	  
	  showFullImage: function(){
		  $('#photoModal').modal()  
	  },
	  
	  initialize : function (){
	  	var params = {
		  	'username' : app.Session.username,
		  	'type' : 'target'
	  	}
	  	this.model = new app.Models.User(params)
		  this.listenTo(this.model, 'change', this.render)
		  this.listenTo(this.model, 'fetch', this.render)
//		  this.listenTo(this.model, 'destroy', )		  		  
	  },
	  kill: function() {
		var secret = this.$el.find('#secret').val();
		var url = WEB_ROOT + 'users/' +  + '/target/'
		var view = this;
		this.model.destroy({
			headers: {'X-DMAssassins-Secret': secret},
			success: function(){
				view.initialize()
				view.render()
			}
			
			})


	  },
	  render: function(){
		this.$el.html( this.template ( this.model.attributes ) );
		return this;  
	  }	    
  })
  
})(jQuery);