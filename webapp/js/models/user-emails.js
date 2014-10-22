//
// js/models/user-emails.js
// dmassassins.js
//
// Copyright (c) 2014 Matt Gerstman
// MIT License.
//
// Rules model, loads rules from the db so that admins can define custom rules per game


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

    app.Models.UserEmails = Backbone.Model.extend({
        defaults: {
            alive: ['Loading...'],
            all: ['Loading...']
        },
        url: function() {
            var game_id = app.Running.Games.getActiveGameId();
            return config.WEB_ROOT + 'game/' + game_id + '/users/email/';
        }
    });
})();