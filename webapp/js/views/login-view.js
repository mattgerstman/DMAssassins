// shows the login screen
// js/views/login-view.js

var app = app || {
    Collections: {},
    Models: {},
    Views: {},
    Routers: {},
    Running: {},
    Session: {}
};

(function($){
 'use strict';
  app.Views.LoginView = Backbone.View.extend({
	   
	     
	  template: _.template( $('#login-template').html() ),
	  
	  events: {
	  	'click .btn-facebook' : 'login'
	  },
	  
	  initialize : function (){	  
	  	//console.log('yo');
		this.model = app.Session;
	  },
	  // call the model login function
	  login: function(){
  		//console.log('login');
  		this.model.login()
	  },
	  // render the login page
	  render: function(){
		this.$el.html( this.template ( this.model.attributes ) );
		return this;  
	  },
  })
  
})(jQuery);