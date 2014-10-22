// Users Collection. Handles all of the users for a game for an admin
// js/collections/users.js
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
    app.Collections.Teams = Backbone.Collection.extend({
		// Collection of teams
        model: app.Models.Team,
        // The api is going to return a mapping, parse to an array
        parse: function(response){
            return _.values(response);
        },
        url: function(){
            var game_id = app.Running.Games.getActiveGameId();
            return config.WEB_ROOT + 'game/' + game_id + '/team/';
        }

    });
})();
