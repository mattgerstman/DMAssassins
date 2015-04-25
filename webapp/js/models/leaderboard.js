//
// js/models/leaderboard.js
// dmassassins.js
//
// Copyright (c) 2014 Matt Gerstman
// MIT License.
//
// Leaderboard model, displays high scores for a game

(function() {
    'use strict';

    app.Models.Leaderboard = Backbone.Model.extend({
        defaults: {
            teams_enabled: true,
            user_col_width: 20,
            team_col_width: 20,
            users: [{
                name: strings.loading,
                kills: strings.loading,
                team_name: strings.loading
            }],
            teams: {
                "Loading...": {
                    alive   : strings.loading,
                    kills   : strings.loading,
                    players : strings.loading
                }
            }
        },
        parse: function(data) {
            data.user_col_width = 100 / (3 + this.get('teams_enabled'));
            data.team_col_width = 20;
            return data;
        },
        fetch: function(options) {
            if (app.Running.Games.getActiveGameId() === null) {
                return;
            }
            return Backbone.Model.prototype.fetch.call(this, options);
        },
        url: function() {
            var game_id = app.Running.Games.getActiveGameId();
            return config.WEB_ROOT + 'game/' + game_id + '/leaderboard/';
        }
    });
})();
