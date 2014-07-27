// js/views/login-view.js

var app = app || {Models:{}, Views:{}, Routers:{}, Running:{}};

(function($){
 'use strict';
  app.Views.LoginView = Backbone.View.extend({
	   
	     
	  template: _.template( $('#login-template').html() ),
	  el: '#login',
	  
	  events: {
	  	'click .btn-facebook' : 'login'
	  },
	  model: new app.Models.Login(),
	  
	  initialize : function (){	  
	  
	  },
	  login: function(){
		  this.model.login()
		  app.Running.Router.navigate('target', true)
	  },
	  render: function(){
		this.$el.html( this.template ( this.model.attributes ) );
		return this;  
	  },
  })
  
})(jQuery);