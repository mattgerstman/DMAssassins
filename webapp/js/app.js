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

    app.Running.appView = new app.Views.AppView();
    app.Running.appView.render();


    app.Session = new app.Models.Session();
    app.Session.setAuthHeader();

    app.Running.Games = new app.Collections.Games();

    var user_id = app.Session.get('user_id');

    app.Running.ProfileModel = new app.Models.Profile(app.Session.get('user'))

    app.Running.TargetModel = new app.Models.Target({
        assassin_id: user_id
    })
    app.Running.LeaderboardModel = new app.Models.Leaderboard();
    app.Running.RulesModel = new app.Models.Rules();

    app.Running.Router = new app.Routers.Router();
    Backbone.history.start();

});