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
            'change #games': 'checkFields'

        },
        // previous page, may depricate
        loaded_from: 'login',
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
            if (!!app.Running.Games.getActiveGameId()) {
                app.Running.Router.back();
                return;
            }
            $('.select-game-active').addClass('hide').removeClass('select-game-active');
            $('#create-or-join').removeClass('hide');
            $('.logo').removeClass('hide');
        },
        // show the create game s ubview
        createGame: function(event) {
            event.preventDefault();
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
        joinGame: function(event) {
            event.preventDefault();            
            var selected = this.$el.find('#games option:selected');
            var game_id = selected.val();
            var need_password = selected.attr('game_has_password') == 'true';
            var password = need_password ? $('#join_game_password').val() : '';
            
            var teams_enabled = this.$el.find('#join_game_team').attr('disabled') != 'disabled';
            var team_id = this.$el.find('#join_game_team option:selected').val();
            
            if (teams_enabled && !team_id) {
                this.badTeam();
                return;
            }            
            app.Running.Games.joinGame(game_id, password, team_id);
        },
        badPassword: function(){
            $('#join_password_block').addClass('has-error');
            $('label[for=join_game_password]').text('Invalid Password:');
        },
        badTeam: function(){
            $('#join_team_block').addClass('has-error');
            $('label[for=join_game_team]').text('Must Select A Team:');
        },
        fixFields: function(){
          this.$el.find('.has-error').removeClass('has-error');
          this.$el.find('label[for=join_game_password]').text('Password:');
          this.$el.find('label[for=join_game_team]').text('Team:');
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
            $('#create_game_password').attr('disabled', !e.target.checked);
        },
        noTeams: function(){
            var teamField = this.$el.find('#join_game_team');
            teamField.attr('disabled', true);            
            teamField.find('#team_placeholder').text('No Teams');
        },
        // toggles the password entry field and team entry fieldon join game
        checkFields: function() {
            // Password field setup
            this.fixFields();
            var selected = this.$el.find('#games option:selected');
            var need_password = selected.attr('game_has_password') == 'true';
            var passwordField = this.$el.find('#join_game_password');
            var passwordPlaceholder = need_password ? '' : 'No Password';
            passwordField.attr('disabled', !need_password);
            passwordField.val(passwordPlaceholder);
            
            var game_id = selected.val();
            var game = app.Running.Games.get(game_id);
            if (!game) {
                return;
            }

            var teamField = this.$el.find('#join_game_team');
            teamField.find('#team_placeholder').text('Loading..');

            var that = this;
            var url = config.WEB_ROOT + 'game/' + game_id + '/team/';
            app.Running.Teams.fetch({
                url:url,
                success:function(teams){        
                    teamField.attr('disabled', false);
                    var teamOptionsTemplate = _.template($('#select-game-team-option').html());
                    var teamOptionsHTML = teamOptionsTemplate({teams: app.Running.Teams.toJSON()});
                    teamField.html(teamOptionsHTML);
                    if (!app.Running.Teams.length){
                        that.noTeams();
                    }
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