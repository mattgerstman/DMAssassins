//
// js/models/team.js
// dmassassins.js
//
// Copyright (c) 2014 Matt Gerstman
// MIT License.
//
// Team model, manages single team

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

    app.Models.Team = Backbone.Model.extend({

        // default profile properties
        defaults: {
            'user_id': null,
            'team_id': null,
            'team_name': ''
        },
        idAttribute : 'team_id',
        fetch: function(options) {
            if (app.Running.Games.getActiveGameId() === null) {
                return;
            }
            return Backbone.Model.prototype.fetch.call(this, options);
        },
        url: function(){
            var game_id = app.Running.Games.getActiveGameId();
            var user_id = this.get('user_id');
            var team_id = this.get('team_id');

            if (!!user_id) {
                return config.WEB_ROOT + 'game/' + game_id + '/user/' + user_id + '/team/' + team_id + '/';
            }
            if (!!team_id) {
                return config.WEB_ROOT + 'game/' + game_id + '/team/' + team_id + '/';
            }
            return config.WEB_ROOT + 'game/' + game_id + '/team/';
        }

    });
})();
