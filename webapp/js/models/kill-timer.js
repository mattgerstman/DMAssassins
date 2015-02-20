//
// js/models/kill-timer.js
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

    app.Models.KillTimer = Backbone.Model.extend({

        // default profile properties
        defaults: {
            'game_id': '',
            'execute_ts': 0,
            'create_ts': 0
        },
        idAttribute: 'game_id',
        url: function(){
            var game_id = app.Running.Games.getActiveGameId();
            return config.WEB_ROOT + 'game/' + game_id + '/kill_timer/';
        }

    });
})();
