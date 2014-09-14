//
// js/models/profile.js
// dmassassins.js
//
// Copyright (c) 2014 Matt Gerstman
// MIT License.
//
// Profile model, manages logged in user

(function() {
    'use strict';

    app.Models.Profile = Backbone.Model.extend({

        // default profile properties
        defaults: {
            'user_id': '',
            'username': '',
            'email': '',
            'properties': {
                'name': 'Loading..',
                'facebook': 'Loading..',
                'secret': 'Loading..',
                'photo_thumb': SPY,
                'photo': SPY
            }

        },
        idAttribute : 'user_id',
        url: function() {
            var game_id = app.Running.Games.getActiveGameId();
            return config.WEB_ROOT + 'game/' + game_id + '/users/' + this.get('user_id') + '/';           
        },       
        joinGame: function(game_id, game_password) {
            var that = this;            
            this.save(null, {
                headers: {
                    'X-DMAssassins-Game-Password': game_password
                },
                success: function() {
                    that.trigger('join-game');
                    Backbone.history.navigate('my_profile', {
                        trigger: true
                    });
                }
            });
        },
        quit: function(secret) {
            var that = this;
            this.destroy({
                headers: {
                    'X-DMAssassins-Secret': secret
                },
                success: function() {
                    if (!app.Running.Games.removeActiveGame()) {
                        Backbone.history.navigate('#logout', {
                            trigger: true
                        });
                        return;
                    }
                }
            })
        },
    })
})();