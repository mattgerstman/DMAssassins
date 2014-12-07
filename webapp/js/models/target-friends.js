//
// js/models/target-friends.js
// dmassassins.js
//
// Copyright (c) 2014 Matt Gerstman
// MIT License.
//
// model for target pages

var app = app || {
    Collections: {},
    Models: {},
    Views: {},
    Routers: {},
    Running: {},
    Session: {}
};

(function() {
    'use strict';
    app.Models.TargetFriends = Backbone.Model.extend({
        defaults: {
            'count': 0,
            'friends': []
        },
        url: function() {
            var game_id = app.Running.Games.getActiveGameId();
            var user_id = app.Running.User.get('user_id');
            return config.WEB_ROOT + "game/" + game_id + '/user/' + user_id + '/target/friends/';
        },
    });
})();
