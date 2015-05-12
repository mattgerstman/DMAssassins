//
// js/models/target-friends.js
// dmassassins.js
//
// Copyright (c) 2014 Matt Gerstman
// MIT License.
//
// model for target pages

(function() {
    'use strict';
    app.Models.TargetFriends = Backbone.Model.extend({
        defaults: {
            'count': 0,
            'friends': []
        },
        parse: function(response) {
            response.friends = response.friends || [];
            return response;
        },
        fetch: function(options) {
            if (app.Running.Games.getActiveGameId() === null) {
                return;
            }
            return Backbone.Model.prototype.fetch.call(this, options);
        },
        url: function() {
            var game_id = app.Running.Games.getActiveGameId();
            var user_id = app.Running.User.get('user_id');
            return config.WEB_ROOT + "game/" + game_id + '/user/' + user_id + '/target/friends/';
        },
    });
})();
