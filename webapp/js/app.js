var app = app || {Models:{}, Views:{}, Routers:{}, Running:{}};
$(function(){
	'use strict';	
	app.Running.Router = new app.Routers.Router();
	Backbone.history.start();
	
});