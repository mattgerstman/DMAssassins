//
// js/models/plot-twist.js
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

    app.Models.PlotTwist = Backbone.Model.extend({

        // default profile properties
        defaults: {
            'plot_twist_name': '',
            'send_email': false
        },
        url: function(){
            var game_id = app.Running.Games.getActiveGameId();
            return config.WEB_ROOT + 'game/' + game_id + '/plot_twist/';
        }

    });
})();
