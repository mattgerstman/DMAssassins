//
// js/models/kill-timer.js
// dmassassins.js
//
// Copyright (c) 2014 Matt Gerstman
// MIT License.
//
// Team model, manages single team

(function() {
    'use strict';

    app.Models.KillTimer = Backbone.Model.extend({

        // default profile properties
        defaults: {
            'game_id': null,
            'execute_ts': 0,
            'create_ts': 0
        },
        idAttribute: 'game_id',
        fetch: function(options) {
            if (this.get('game_id') === null) {
                return;
            }
            return Backbone.Model.prototype.fetch.call(this, options);
        },
        url: function(){
            var game_id = app.Running.Games.getActiveGameId();
            return config.WEB_ROOT + 'game/' + game_id + '/kill_timer/';
        }

    });
})();
