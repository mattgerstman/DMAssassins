var app = app || {Models:{}, Views:{}, Routers:{}, Running:{}, Session:{}};

$(function(){
	'use strict';	
	app.Running.appView = new app.Views.AppView();
	app.Running.appView.render();	

	app.Session = new app.Models.Session();

	app.Running.UserModel = new app.Models.Profile(app.Session.get('user'))
	app.Running.TargetModel = new app.Models.Target(app.Session.get('target'))

	app.Running.Router = new app.Routers.Router();
	Backbone.history.start();
	
});