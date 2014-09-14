// Select Games Model, Handles game creation, selection, and joining
// Focuses on game mappings on the server
// js/models/games.js
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
        // handle on initiliazation
        url:function(){
            var user_id = app.Session.get('user_id');
            return config.WEB_ROOT + 'users/' + user_id + '/game/';
            
        },
        initialize: function() {
            
            var sessionGame = app.Session.get('game');
            if (sessionGame)
            {
                this.active_game = new app.Models.Game(sessionGame);
                return;
            }
            
            this.active_game = new app.Models.Game();
            
            
        },
        addGame: function(game_id){
            var that = this;
            var game = new app.Models.Game({game_id: game_id});
            game.fetch({success: function(yo){
                console.log(yo);
                that.add(game);
                that.setActiveGame(game.get('game_id'));
            }})
        },
        removeActiveGame: function(){
            this.remove(this.active_game);
            var newGame = this.findWhere({game_started: true})
            if (!newGame)
            {
                newGame = this.findWhere({game_started: false})
            }
            this.setActiveGame(newGame);
            return newGame;
        },
        setActiveGame: function(game_id) {
            var game = this.get(game_id);
            if (!game) {
                return null;
            }
            this.active_game = game;
            this.trigger('game-change');
            app.Session.set('game', JSON.stringify(this.active_game));
            return this.active_game;
        },
        getActiveGame: function() {        
            return this.active_game;

        },
        getActiveGameId: function() {        
            return this.active_game.get('game_id');
        }

        
        
    })
})();
