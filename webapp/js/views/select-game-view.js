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
        // shows the create game subview
        showCreateGame: function() {
            $('.js-logo').addClass('hide');
            $('.js-create-or-join').addClass('hide');
            $('.js-create-game').addClass('js-select-game-active');
            $('.js-create-game').removeClass('hide');
        },
        // shows the join game subview
        showJoinGame: function() {
            $('.js-logo').addClass('hide');
            $('.js-create-or-join').addClass('hide');
            $('.js-join-game').addClass('js-select-game-active');
            $('.js-join-game').removeClass('hide');
        },
        // cancels the game creation/selection
        goBack: function() {
            if (!!app.Running.Games.getActiveGameId()) {
                app.Running.Router.back();
                return;
            }
            $('.js-select-game-active').addClass('hide').removeClass('js-select-game-active');
            $('.js-create-or-join').removeClass('hide');
            $('.js-logo').removeClass('hide');
        },
        // show the create game s ubview
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
        joinGame: function(event) {
            event.preventDefault();
            var selected = this.$el.find('#js-select-game option:selected');
            var game_id = selected.val();
            var need_password = selected.data('game-has-password');
            var password = need_password ? $('#js-join-game-password').val() : '';

            var teams_enabled = this.$el.find('#js-join-game-team').attr('disabled') != 'disabled';
            var team_id = this.$el.find('#js-join-game-team option:selected').val();

            if (teams_enabled && !team_id) {
                this.badTeam();
                return;
            }
            $(event.currentTarget).addClass('disabled').text('Joining');
            app.Running.Games.joinGame(game_id, password, team_id);
        },
        badPassword: function(){
            $('#join_password_block').addClass('has-error');
            $('label[for=js-join-game-password]').text('Invalid Password:');
            $('.js-join-game-submit').removeClass('disabled').text('Join');
        },
        badTeam: function(){
            $('#join_team_block').addClass('has-error');
            $('label[for=js-join-game-team]').text('Must Select A Team:');
        },
        fixFields: function(){
            this.$el.find('.has-error').removeClass('has-error');
            this.$el.find('label[for=js-join-game-password]').text('Password:');
            this.$el.find('label[for=js-join-game-team]').text('Team:');
        },
        // finish up and navigate to your profile
        finish: function(game) {
            app.Running.Games.setActiveGame(game.get('game_id'));
            Backbone.history.navigate('my_profile', {
                trigger: true
            });
        },
        finishJoin: function(){
            Backbone.history.navigate('my_profile', {
                trigger: true
            });
        },
        // toggles the password entry field on create game
        togglePassword: function(e) {
            $('#js-create-game-password').attr('disabled', !e.target.checked);
        },
        noTeams: function(){
            var teamField = this.$el.find('#js-join-game-team');
            teamField.attr('disabled', true);
            teamField.find('#js-team-placeholder').text('No Teams');
        },
        // toggles the password entry field and team entry field on join game
        checkFields: function() {
            // Password field setup
            this.fixFields();
            // get selected game
            var selected = this.$el.find('#js-select-game option:selected');
            // check if it needs a password
            var need_password = selected.data('game-has-password');

            // grab the password field
            var passwordField = this.$el.find('#js-join-game-password');

            // Set it to no password if we don't have one
            var passwordPlaceholder = need_password ? '' : 'No Password';

            // If we do have one mark it as not disabled
            passwordField.attr('disabled', !need_password);
            passwordField.val(passwordPlaceholder);

            // Get game model
            var game_id = selected.val();
            var game = app.Running.Games.get(game_id);
            if (!game) {
                return;
            }

            // Set teams to loading if the game has teams
            var teamField = this.$el.find('#js-join-game-team');
            teamField.find('#js-team-placeholder').text('Loading..');

            // Set the teams
            var that = this;
            var url = config.WEB_ROOT + 'game/' + game_id + '/team/';
            app.Running.Teams.fetch({
                url:url,
                success:function(teams){
                    // if there are no teams handle that appropraitely
                    if (!app.Running.Teams.length){
                        that.noTeams();
                    }
                    // if we have teams mark it as not disabled
                    teamField.attr('disabled', false);
                    // render team select
                    var teamOptionsTemplate = _.template($('#select-game-team-option').html());
                    var teamOptionsHTML = teamOptionsTemplate({teams: app.Running.Teams.toJSON()});
                    teamField.html(teamOptionsHTML);
                },
                error:function(){
                    that.noTeams();
                }
            });
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
})(jQuery);
