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
            game_has_password: false,
            member: true
        },

        idAttribute: 'game_id',
        urlRoot: config.WEB_ROOT + 'users/',
        url: function() {
            var url = this.urlRoot;
            url += app.Session.get('user_id') + '/game/';            
            var game_id = this.get('game_id');
            if (!game_id)
            {
                return url;
            }            
            return url + game_id + '/';
        }
    })
})();