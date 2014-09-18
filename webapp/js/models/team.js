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
            'team_id': '',
            'team_name': ''
        },
        url: function(){            
            var game_id = app.Running.Games.getActiveGameId();
            return config.WEB_ROOT + 'game/' + game_id + '/team/' + this.get('team_id') + '/';           
        }

    })
})();
