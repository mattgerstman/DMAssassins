//
// js/views/nav-game-view.js
// dmassassins.js
//
// Copyright (c) 2014 Matt Gerstman
// MIT License.
//
// handles the game dropdown in the nav


var app = app || {
    Collections: {},
    Models: {},
    Views: {},
    Routers: {},
    Running: {},
    Session: {}
};

(function($) {
    'use strict';
    app.Views.NavGameView = Backbone.View.extend({
        template: _.template($('#nav-game-template').html()),
        el: '#games_dropdown',

        tagName: 'ul',

        events: {
            'click li a.switch_game': 'select'
        },
        // constructor, loads a user id so we can get their games from the model
        initialize: function() {
            this.collection = app.Running.Games;
            this.listenTo(this.collection, 'fetch', this.render);
            this.listenTo(this.collection, 'change', this.render);
            this.listenTo(this.collection, 'reset', this.render);
            this.listenTo(this.collection, 'add', this.render);
            this.listenTo(this.collection, 'remove', this.render);
            this.listenTo(this.collection, 'game-change', this.render);

        },
        handleJoin: function () {
            var availableGame = _.findWhere(this.collection.toJSON(), {member: false});
            if (availableGame === undefined)
            {
                this.hideJoin();
                return;
            }
            this.showJoin();
        },
        hideJoin: function () {
            this.$el.find('#nav_join_game').addClass('hide');
        },
        showJoin: function () {
            this.$el.find('#nav_join_game').removeClass('hide');
        },
        showCurrentGame: function() {
            var game_id = app.Running.Games.getActiveGameId();
            this.$el.find('#nav_' + game_id).removeClass('hide');
        },
        updateText: function() {

            $('.game_name').removeClass('hide');
            if (Backbone.history.fragment == 'join_game') {
                this.showCurrentGame();
                $('#games_header').text('Join Game');
                $('#games_header_short').text('Join Game');
                return this;
            }

            if (Backbone.history.fragment == 'create_game') {
                this.showCurrentGame();
                $('#games_header').text('Create Game');
                $('#games_header_short').text('Create Game');
                return this;
            }

            var game = this.collection.getActiveGame();
            if (!game)
            {
                return this;
            }
            var game_name = game.get('game_name');
            $('#games_header').text(game_name);
            var max = 9;
            if (game_name.length > max) {
                game_name = game_name.substr(0, max - 3) + '...';
            }
            $('#games_header_short').text(game_name);

            var game_id = app.Running.Games.getActiveGameId();
            this.$el.find('#nav_' + game_id).addClass('hide');
        },
        // loads the items into the dropdown and changes the dropdown title to the current game
        render: function() {
            this.$el.html(this.template({
                games: _.where(this.collection.toJSON(), {
                    member: true
                })
            }));
            this.handleJoin();
            this.updateText();
            return this;

        },
        // select a game from the dropdown
        select: function(event) {
            var game_id = $(event.target).attr('game_id');
            app.Running.Games.setActiveGame(game_id);

        }
    });

})(jQuery);