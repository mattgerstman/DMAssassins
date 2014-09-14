//
// js/models/games.js
// dmassassins.js
//
// Copyright (c) 2014 Matt Gerstman
// MIT License.
//
// Single game model

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
    app.Models.Game = Backbone.Model.extend({

        // default properties with a fake game
        defaults: {
            game_id: '',
            game_name: '',
            game_started: false,
            game_has_password: false
        },

        idAttribute: 'game_id',
        urlRoot: config.WEB_ROOT + 'game/',
        url: function() {
            return this.urlRoot + this.get('game_id') + '/';
        }
    })
})();