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

        model: app.Models.Game,

        parse: function(response) {
            if (!response) {
                return null;
            }
            var list = [];
            _.each(response.available, function(item){
                item.member = false;
                list.push(item);
            })
            _.each(response.member, function(item){
                item.member = true;
                list.push(item);
            })
            return list;
        },
        comparator: 'game_name',
        active_game: null,
        // handle on initiliazation
        url:function(){
            var user_id = app.Session.get('user_id');
            return config.WEB_ROOT + 'user/' + user_id + '/game/';
            
        },
        addGame: function(game_id){
            var that = this;
            var game = new app.Models.Game({game_id: game_id});
            game.fetch({success: function(){
                that.add(game);
                that.setActiveGame(game.get('game_id'));
            }})
        },
        joinGame: function(game_id, password, team_id) {
            app.Running.User.joinGame(game_id, password, team_id);  
            this.trigger('game-change');          
        },
        setArbitraryActiveGame: function(silent) {
            var newGame = this.findWhere({game_started: true})
            if (!newGame)
            {
                newGame = this.findWhere({game_started: false})
            }
            this.setActiveGame(newGame, silent);
            return newGame;
        },
        removeActiveGame: function(){
            this.remove(this.active_game);
            return this.setArbitraryActiveGame();
        },
        setActiveGame: function(game_id, silent) {
            var game = this.get(game_id);
            if (!game) {
                return null;
            }
            game.fetchProperties();
            this.active_game = game;
            app.Session.set('game_id', game_id);
            app.Session.set('has_game', true);

            if (silent === undefined || !silent)
            {   
                this.trigger('game-change');    
            }
            
            return this.active_game;
        },
        getActiveGame: function() {    
            if (!this.active_game)
            {
                this.setArbitraryActiveGame(true);
            }    
            return this.active_game;

        },
        getActiveGameId: function() { 
            var game = this.getActiveGame();
            if (!game)
            {
                return null;
            }
            return game.get('game_id');
        }

        
        
    })
})();
