var app = app || {Models:{}, Views:{}, Routers:{}, Running:{}, Session:{}};

$(function(){
	'use strict';	
	app.Running.appView = new app.Views.AppView();
	app.Running.appView.render();	


	app.Session = new app.Models.Session();
	app.Session.setAuthHeader();
	
	var user_id = app.Session.get('user_id');
	app.Running.GamesModel = new app.Models.Games({user_id:user_id});

	app.Running.NavModel = new app.Models.Nav();
	app.Running.ProfileModel = new app.Models.Profile(app.Session.get('user'))
	
	app.Running.TargetModel = new app.Models.Target(app.Session.get('target'))
	app.Running.LeaderboardModel = new app.Models.Leaderboard();
	app.Running.RulesModel = new app.Models.Rules();
	
	app.Running.Router = new app.Routers.Router();
	Backbone.history.start();
	
});