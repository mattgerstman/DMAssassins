var app = app || {Models:{}, Views:{}, Routers:{}, Running:{}, Session:{}};
$(function(){
	'use strict';	
	app.Running.Router = new app.Routers.Router();
	Backbone.history.start();
	
});