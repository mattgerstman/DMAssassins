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
    app.Views.SelectGameView = Backbone.View.extend({


        template: app.Templates["select-game"],
        tagName: 'div',
        events: {
            'click  .js-create-game-submit'          : 'createGame',
            'click  #js-create-game-need-password'   : 'togglePassword',
            'click  .js-create-or-join-back'         : 'goBack',
            'click  .js-join-game-submit'            : 'joinGame',
            'change #js-select-game'                 : 'checkFields',
            'click  .js-show-create-game'            : 'showCreateGame',
            'click  .js-show-join-game'              : 'showJoinGame'
        },
        // constructor
        initialize: function() {
            this.collection = app.Running.Games;
            this.listenTo(this.collection, 'reset', this.render);
            this.listenTo(this.collection, 'fetch', this.render);
            this.listenTo(app.Running.User, 'join-error-password', this.badPassword);
            this.listenTo(app.Running.User, 'join-game', this.finishJoin);
            this.listenTo(app.Running.Games, 'game-change', this.finishJoin);
        },
        // cancels the game creation/selection
        goBack: function() {
            app.Running.Router.back();
        },
        // show the create games ubview
        createGame: function(event) {
            event.preventDefault();
            var name = $('#create_game_name').val();
            var password = null;
            if ($('#js-create-game-need-password').is(':checked')) {
                password = $('#js-create-game-password').val();

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
            Backbone.history.navigate('my_profile', {
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
            this.checkFields();
            return this;
        }

    });
})();
