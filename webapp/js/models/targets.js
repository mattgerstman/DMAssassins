// Targets Collection. Handles all of the targets for a game for a super admin
// js/collections/targets.js
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
    app.Models.Targets = Backbone.Model.extend({
        defaults : {
            loops: [
                [
                    {
                        "assassin_id": "Loading...",
                        "assassin_name": "Loading...",
                        "assassin_team_name": "Loading...",
                        "assassin_user_role": "Loading...",
                        "assassin_kills": "Loading...",
                        "target_id": "Loading...",
                        "target_name": "Loading...",
                        "target_team_name": "Loading...",
                        "target_user_role": "Loading...",
                        "target_kills": "Loading..."
                    }
                ]
            ]
        },
        // The api is going to return a mapping, parse to an array
        parse: function(response) {
            var targets = response;            
            var data = {loops:[]};
            
            // don't stop until we've hit every user
            while (!_.isEmpty(targets)) {
                // Create a loop
                var loop = [];
                var targetIds = _.keys(targets);
                var firstId = targetIds[0];
                var nextId = null;
                
                // push the first target into the loop
                loop.push(targets[firstId]);
                nextId = targets[firstId].target_id;
                
                // remove the first target from the loop
                delete targets[firstId];
                
                // inner loop until we get back to the first target
                while (nextId != firstId) {
                    var currentTarget = targets[nextId];
                    // push the next target to the loop
                    loop.push(currentTarget);
                    // delete the last id
                    delete targets[nextId];
                    // move on to the next id
                    nextId = currentTarget.target_id;


                }
                // push this loop on the loops container
                data.loops.push(loop);
            }
            return data;
        },
        url: function(){
            var game_id = app.Running.Games.getActiveGameId();
            return config.WEB_ROOT + 'game/' + game_id + '/targets/';
        }

    });
})();
