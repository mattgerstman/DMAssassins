//
// js/models/target.js
// dmassassins.js
//
// Copyright (c) 2014 Matt Gerstman
// MIT License.
//
// model for target pages

(function() {
    'use strict';
    app.Models.Target = Backbone.Model.extend({
        defaults: {
            'game_id': null,
            'assassin_id': null,
            'facebook_id': null,
            'username': null,
            'user_id': null,
            'properties': {
                'name': strings.loading,
                'facebook': strings.loading,
                'team':strings.loading,
                'photo_thumb': SPY,
                'photo': SPY
            }
        },
        url: function() {
            var game_id = app.Running.Games.getActiveGameId();
            return config.WEB_ROOT + "game/" + game_id + '/user/' + app.Session.get('user_id') + '/target/';
        },
        idAttribute: 'user_id',
        // constructor
        initialize: function() {
            var target_id = app.Session.get('target_id');
            this.set('user_id', target_id);
            this.listenTo(this, 'fetch', this.saveUserId);
            this.listenTo(this, 'change', this.saveUserId);
            this.listenTo(this, 'reset', this.saveUserId);
        },
        saveUserId: function() {
            var user_id = this.get('user_id');
            app.Session.set('target_id', user_id);
        },
        fetch: function(options) {
            if (app.Running.Games.getActiveGameId() === null) {
                return;
            }
            return Backbone.Model.prototype.fetch.call(this, options);
        }
    });
})();
