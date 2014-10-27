//
// app.js
// dmassassins.js
//
// Copyright (c) 2014 Matt Gerstman
// MIT License.
//

var app = app || {
    Collections: {},
    Models: {},
    Views: {},
    Routers: {},
    Running: {},
    Session: {}
};

// Instantiates all of the running models, routers, and session

$(function() {
    'use strict';

    Raven.config(config.SENTRY_DSN, {
    }).install();

    app.Running.AppView = new app.Views.AppView();
    app.Running.AppView.render();

    app.Session = new app.Models.Session();
    app.Session.setAuthHeader();

    app.Running.Games = new app.Collections.Games();
    app.Running.User = new app.Models.User();
    app.Running.TargetModel = new app.Models.Target();
    app.Running.LeaderboardModel = new app.Models.Leaderboard();
    app.Running.RulesModel = new app.Models.Rules();

    app.Running.Users = new app.Collections.Users();
    app.Running.Teams = new app.Collections.Teams();

    app.Running.User.listenTo(app.Running.Games, 'game-change', app.Running.User.fetch);
    app.Running.TargetModel.listenTo(app.Running.Games, 'game-change', app.Running.TargetModel.fetch);
    app.Running.LeaderboardModel.listenTo(app.Running.Games, 'game-change', app.Running.LeaderboardModel.fetch);
    app.Running.RulesModel.listenTo(app.Running.Games, 'game-change', app.Running.RulesModel.fetch);

    app.Running.User.listenTo(app.Running.User, 'fetch', app.Running.User.checkAccess);
    app.Running.User.listenTo(app.Running.User, 'change', app.Running.User.checkAccess);

    app.Running.Router = new app.Routers.Router();
    Backbone.history.start();

});