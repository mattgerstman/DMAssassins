//
// js/models/leaderboard.js
// dmassassins.js
//
// Copyright (c) 2014 Matt Gerstman
// MIT License.
//
// Leaderboard model, displays high scores for a game

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

    app.Models.Leaderboard = Backbone.Model.extend({
        defaults: {
            teams_enabled: true,
            user_col_width: 20,
            team_col_width: 20,
            users: [{
                name: "Loading...",
                kills: "Loading...",
                team_name: "Loading..."
            }, {
                name: "Loading...",
                kills: "Loading...",
                team_name: "Loading..."
            }],
            teams: [{
                "Loading...": "Loading..."
            }]
        },
        parse: function(data) {
            data.user_col_width = 100 / (3 + this.get('teams_enabled'));
            data.team_col_width = 20;
            return data;
        },
        url: function() {
            var game_id = app.Running.Games.getActiveGameId();
            if (!game_id)
                return null;
            return config.WEB_ROOT + 'game/' + game_id + '/leaderboard/';
        }
    });
})();