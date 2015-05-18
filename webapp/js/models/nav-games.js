//
// js/models/nav-games.js
// dmassassins.js
//
// Copyright (c) 2014 Matt Gerstman
// MIT License.
//
// model for nav

(function() {
    'use strict';
    app.Models.NavGames = Backbone.Model.extend({
        defaults: {
            short_text: strings.loading,
            long_text: strings.loading,
            active_game_id: null,
            show_join: true
        },
        initialize: function() {
            this.listenTo(app.Running.Games, 'game-change', this.updateData);
            this.listenTo(app.Running.Games, 'change', this.updateData);
            this.listenTo(app.Running.Games, 'add', this.updateData);
            this.listenTo(app.Running.Games, 'remove', this.updateData);
            this.listenTo(app.Running.Games, 'reset', this.updateData);
            this.listenTo(app.Running.Games, 'fetch', this.updateData);
            this.updateData();
        },
        canJoinGame: function() {
            var availableGame = _.findWhere(app.Running.Games.toJSON(), {member: false});
            if (availableGame === undefined)
            {
                return false;
            }
            return true;
        },
        getGames: function() {
            return _.where(app.Running.Games.toJSON(), {
                    member: true
            });
        },
        // determines what to show in the top bar
        updateData: function() {
            var show_join = this.canJoinGame();
            var games = this.getGames();
            var data = {
                show_join: show_join,
                games: games
            };

            if (Backbone.history.fragment === 'join-game') {
                data.short_text     = strings.join_game;
                data.long_text      = strings.join_game;
                data.active_game_id = 'join_game';
                return this.set(data);
            }

            if (Backbone.history.fragment === 'create-game') {
                data.short_text     = strings.join_game;
                data.long_text      = strings.join_game;
                data.active_game_id = 'create_game';
                return this.set(data);
            }

            var game_name = app.Running.Games.getActiveGameName();
            var short_game_name = game_name;

            var max = 9;
            if (game_name.length > max) {
                short_game_name = game_name.substr(0, max - 3) + '...';
            }

            data.long_text      = game_name;
            data.short_text     = short_game_name;
            data.active_game_id = app.Running.Games.getActiveGameId();
            return this.set(data);
        },
    });
})();
