// Targets Collection. Handles all of the targets for a game for a super admin
// js/collections/targets.js

(function() {
    'use strict';
    app.Models.Targets = Backbone.Model.extend({
        defaults : {
            loops: [
                [
                    {
                        "assassin_id": strings.loading,
                        "assassin_name": strings.loading,
                        "assassin_team_name": strings.loading,
                        "assassin_user_role": strings.loading,
                        "assassin_kills": strings.loading,
                        "target_id": strings.loading,
                        "target_name": strings.loading,
                        "target_team_name": strings.loading,
                        "target_user_role": strings.loading,
                        "target_kills": strings.loading
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
        fetch: function(options) {
            if (app.Running.Games.getActiveGameId() === null) {
                return;
            }
            return Backbone.Model.prototype.fetch.call(this, options);
        },
        url: function(){
            var game_id = app.Running.Games.getActiveGameId();
            if (!game_id)
                return null;
            return config.WEB_ROOT + 'game/' + game_id + '/targets/';
        }

    });
})();
