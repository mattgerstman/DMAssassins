//
// js/views/select-game-view.js
// dmassassins.js
//
// Copyright (c) 2014 Matt Gerstman
// MIT License.
//
// handles game selection


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
    app.Views.SelectGameView = Backbone.View.extend({


        template: _.template($('#select-game-template').html()),
        tagName: 'div',
        events: {
            'click .show-create-game': 'showCreateGame',
            'click .show-join-game': 'showJoinGame',
            'click .create-game-submit': 'createGame',
            'click .join-game-submit': 'joinGame',
            'click .create-or-join-back': 'goBack',
            'click #create_game_need_password': 'togglePassword',
            'change #games': 'checkPassword'

        },
        // previous page, may depricate
        loaded_from: 'login',
        // constructor
        initialize: function(params) {
            this.collection = app.Running.Games;
            this.listenTo(this.collection, 'change', this.render)
            this.listenTo(this.collection, 'fetch', this.render)
        },
        // shows the create game subview
        showCreateGame: function() {
            $('.logo').addClass('hide');
            $('#create-or-join').addClass('hide');
            $('#create-game').addClass('select-game-active');
            $('#create-game').removeClass('hide');
        },
        // shows the join game subview
        showJoinGame: function() {
            $('.logo').addClass('hide');
            $('#create-or-join').addClass('hide');
            $('#join-game').addClass('select-game-active');
            $('#join-game').removeClass('hide');
        },
        // cancels the game creation/selection
        goBack: function() {
            if (app.Session.get('has_game') == "true") {
                history.back();
                return;
            }
            $('.select-game-active').addClass('hide').removeClass('select-game-active');
            $('#create-or-join').removeClass('hide');
            $('.logo').removeClass('hide');
        },
        // show the create game s ubview
        createGame: function() {
            var name = $('#create_game_name').val();
            var password = null;
            if ($('#create_game_need_password').is(':checked')) {
                password = $('#create_game_password').val();

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
        // loads the join game later view
        loadJoinGame: function(user_id) {
            var that = this;
            that.showJoinGame();
            this.collection.fetch({
                wait: true,
                success: function() {
                    that.showJoinGame();
                }
            });

        },
        // posts to the join game model
        joinGame: function() {
            var game_id = $('#games option:selected').val();
            var password = $('#join_game_password').val();
            app.Running.Games.joinGame(game_id, password);
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
            $('#create_game_password').attr('disabled', !e.target.checked);
        },
        // toggles the password entry field on join game
        checkPassword: function(e) {
            var need_password = $('#games option:selected').attr('game_has_password') == 'true';
            $('#join_game_password').attr('disabled', !need_password);

        },
        render: function() {
            this.$el.html(this.template({
                games: _.where(this.collection.toJSON(), {
                    member: false
                })
            }));
            return this;
        }

    })

})(jQuery);