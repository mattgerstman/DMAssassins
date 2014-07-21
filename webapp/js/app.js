var app = app || {};


$(function(){
	'use strict';	
	app.runningApp = new app.AppView();
	var navView = new app.NavView();
	navView.render();
	var router = new app.Router();
	Backbone.history.start();
 
	
});