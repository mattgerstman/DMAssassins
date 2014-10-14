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

        // default properties for a game to appear "Loading..."
        defaults: {
            game_id: '',
            game_name: 'Loading...',
            game_started: false,
            game_has_password: false,
            member: true
        },

		// sets game_id as the id attribute
        idAttribute: 'game_id',
        // Sets the url root
        urlRoot: config.WEB_ROOT + 'user/',
        // Determine if teams are enabled for this game
        areTeamsEnabled: function() {
            return this.getProperty('teams_enabled') == 'true';
        },
        getProperty: function(key) {
            var properties = this.get('game_properties');
            if (!properties)
                return null;
            if (properties[key] === undefined)
                return null;
            return properties[key];
        },
        fetchProperties: function() {
            var url = this.gameUrl();
            return this.fetch({url: url});
        },
        gameUrl: function() {
            var url = config.WEB_ROOT + 'game/' + this.get('game_id') + '/';
            return url;
        },
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
    });
})();
