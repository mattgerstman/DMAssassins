//
// js/models/target.js
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
    app.Models.Target = Backbone.Model.extend({
        defaults: {
            'game_id': null,
            'assassin_id': '',
            'username': '',
            'user_id': '',
            'properties': {
                'name': 'Loading...',
                'facebook': 'Loading...',
                'team':'Loading...',
                'photo_thumb': SPY,
                'photo': SPY
            }
        },
        url: function() {
            var game_id = app.Running.Games.getActiveGameId();
            return config.WEB_ROOT + "game/" + game_id + '/user/' + this.get('assassin_id') + '/target/';
        },
        // consstructor
        initialize: function() {
            if (!this.get('assassin_id')) {
                this.assassin_id = app.Session.get('user_id');
            }
            this.idAttribute = 'assassin_id'
        }
    })
})();
