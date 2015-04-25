//
// js/models/user.js
// dmassassins.js
//
// Copyright (c) 2014 Matt Gerstman
// MIT License.
//
// Rules model, loads rules from the db so that admins can define custom rules per game

(function() {
    'use strict';

    app.Models.Rules = Backbone.Model.extend({
        defaults: {
            rules: strings.loading
        },
        save: function (attrs, options) {
            options = options || {};
            options.type = 'PUT';
            return Backbone.Model.prototype.save.call(this, attrs, options);
        },
        fetch: function(options) {
            if (app.Running.Games.getActiveGameId() === null) {
                return;
            }
            return Backbone.Model.prototype.fetch.call(this, options);
        },
        url: function() {
            var game_id = app.Running.Games.getActiveGameId();
            return config.WEB_ROOT + 'game/' + game_id + '/rules/';
        }
    });
})();
