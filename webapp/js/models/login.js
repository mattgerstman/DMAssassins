// js/models/login.js

var app = app || {Models:{}, Views:{}, Routers:{}, Running:{}};
(function() {
	'use strict';
	
	app.Models.Login = Backbone.Model.extend({
		defaults: {
				logo: '/assets/img/logo.png'
			},
		login: function(){

			FB.getLoginStatus(function(response){
				  if (response.status === 'connected') {
				    // Logged into your app and Facebook.
						console.log(response);


				  } else if (response.status === 'not_authorized') {
				    // The person is logged into Facebook, but not your app.
				  } else {
				    // The person is not logged into Facebook, so we're not sure if
				    // they are logged into this app or not.
				  }
			})

		}
			

	})
})();