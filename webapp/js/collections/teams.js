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

        model: app.Models.Team,
        parse: function(response){
            return _.values(response);  
        },
        url: function(){
            var game_id = app.Running.Games.getActiveGameId();
            return config.WEB_ROOT + 'game/' + game_id + '/team/';
        }
        
    })
})();
