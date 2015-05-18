// Games Collection, Handles all of the games a user has access to
// js/collections/games.js
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
    app.Collections.Games = Backbone.Collection.extend({

        // Collection of type Model
        model: app.Models.Game,

        // Sort games into those the user is a member of and those they can join
        parse: function(response) {
            if (!response) {
                return null;
            }
            var list = [];
            _.each(response.available, function(item){
                item.member = false;
                list.push(item);
            });
            _.each(response.member, function(item){
                item.member = true;
                list.push(item);
            });
            return list;
        },
        // Used to sort games alphabetically
        comparator: 'game_name',
        // the game the user is currently viewing
        active_game: null,
        // handle on initiliazation
        url:function(){
            return config.WEB_ROOT + 'game/';
        },
        // Add a game to the collection (usually after creation)
        addGame: function(game_id){
            var that = this;
            var game = new app.Models.Game({game_id: game_id});
            game.fetch({success: function(){
                that.add(game);
                that.setActiveGame(game.get('game_id'));
            }});
        },
        // Have a user join a new game
        joinGame: function(game_id, password, team_id) {
            app.Running.User.joinGame(game_id, password, team_id);
        },
        clearActiveGame: function(options) {
            this.active_game = null;
            this.reset(options);
        },
        // Set the active game to an arbitrary one
        setArbitraryActiveGame: function(silent) {
            var newGame = this.findWhere({game_started: true, member:true});
            if (!newGame)
            {
                newGame = this.findWhere({game_started: false, member:true});
            }

            if (!newGame) {
                return null;
            }

            if (!newGame.get('game_id')) {
                return null;
            }

            this.setActiveGame(newGame, silent);
            return newGame;
        },
        // Removes the active game from the collection
        removeActiveGame: function(){
            this.remove(this.active_game);
            return this.setArbitraryActiveGame();
        },
        // Sets the active game id
        setActiveGame: function(game_id, silent) {
            // validate the game_Id
            var game = this.get(game_id);
            if (!game) {
                return null;
            }

            game.fetchProperties();
            this.active_game = game;
            app.Session.set('game_id', game_id);
            app.Session.set('has_game', true);
            app.Session.set('game_started', game.get('game_started'));

            // trigger a game change if necessary
            if (silent === undefined || !silent)
            {
                this.trigger('game-change');
            }
            return this.active_game;
        },

        // Return the active game
        getActiveGame: function() {
            if (!this.active_game)
            {
                this.setArbitraryActiveGame(true);
            }
            return this.active_game;

        },
        // Returns the active game's id
        getActiveGameId: function() {
            var game = this.getActiveGame();
            if (!game)
            {
                return null;
            }
            return game.get('game_id');
        },
        getActiveGameName: function() {
            var game = this.getActiveGame();
            if (!game)
            {
                return strings.loading;
            }
            return game.get('game_name');
        },
        getActiveGameTeamsEnabled: function() {
            var game = this.getActiveGame();
            if (!game)
            {
                return false;
            }
            return game.areTeamsEnabled();

        },
        hasActiveGameStarted: function() {
            var game = this.getActiveGame();
            if (!game)
            {
                return app.Session.get('game_started');
            }
            return game.get('game_started');
        }
    });
})();
