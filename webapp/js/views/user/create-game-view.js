//
// js/views/create-game-view.js
// dmassassins.js
//
// Copyright (c) 2014 Matt Gerstman
// MIT License.
//
// handles game selection

(function() {
    'use strict';
    app.Views.CreateGameView = Backbone.View.extend({
        template: app.Templates["create-game"],
        tagName: 'div',
        events: {
            'click  .js-create-game-submit'          : 'createGame',
            'click  #js-create-game-need-password'   : 'togglePassword',
            'click  .js-go-back'                     : 'goBack',
            'click  .js-join-game-submit'            : 'joinGame',
        },
        // constructor
        initialize: function() {
            this.collection = app.Running.Games;
            this.listenTo(this.collection, 'reset', this.render);
            this.listenTo(this.collection, 'fetch', this.render);
        },
        // cancels the game creation
        goBack: function(e) {
            e.preventDefault();
            app.Running.Router.back();
        },
        // show the create games ubview
        createGame: function(event) {
            event.preventDefault();
            var name = this.$('#js-create-game-name').val();
            var password = null;
            if ($('#js-create-game-need-password').is(':checked')) {
                    password = this.$('#js-create-game-password').val();

            }
            var that = this;
            this.collection.create({
                game_name: name,
                game_password: password
            }, {
                success: function(game) {
                    that.finish(game);
                }
            });
        },
        // finish up and navigate to your profile
        finish: function(game) {
            app.Running.Games.setActiveGame(game.get('game_id'));
            Backbone.history.navigate('my-profile', {
                trigger: true
            });
        },
        // toggles the password entry field on create game
        togglePassword: function(e) {
            this.$('#js-create-game-password').attr('disabled', !e.target.checked);
        },
        render: function() {
            this.$el.html(this.template({
                games: _.where(this.collection.toJSON(), {
                    member: false
                })
            }));
            return this;
        }

    });
})();
