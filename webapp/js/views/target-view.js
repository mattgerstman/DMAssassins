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
	  	this.model = app.Running.TargetModel;
	  	
		this.listenTo(this.model, 'change', this.render)
		this.listenTo(this.model, 'fetch', this.render)
	  },
	  kill: function() {
		var secret = this.$el.find('#secret').val();
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