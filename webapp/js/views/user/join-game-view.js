//
// js/views/join-game-view.js
// dmassassins.js
//
// Copyright (c) 2014 Matt Gerstman
// MIT License.
//
// handles game selection

(function() {
    'use strict';
    app.Views.JoinGameView = Backbone.View.extend({


        template: app.Templates["join-game"],
        tagName: 'div',
        events: {
            'click  .js-go-back'                     : 'goBack',
            'click  .js-join-game-submit'            : 'joinGame',
            'change #js-select-game'                 : 'checkFields',
            'click  .js-show-join-game'              : 'showJoinGame'
        },
        // constructor
        initialize: function() {
            this.collection = app.Running.Games;
            this.collection.fetch({reset:true});
            this.listenTo(this.collection, 'reset', this.render);
            this.listenTo(this.collection, 'fetch', this.render);
            this.listenTo(app.Running.User, 'join-error-password', this.badPassword);
            this.listenTo(app.Running.User, 'join-game', this.finishJoin);
            this.listenTo(app.Running.Games, 'game-change', this.finishJoin);
        },
        // cancels the game selection
        goBack: function(e) {
            e.preventDefault();
            app.Running.Router.back();
        },
        // posts to the join game model
        joinGame: function(event) {
            event.preventDefault();
            var selected = this.$('#js-select-game option:selected');
            var game_id = selected.val();

            var game = this.collection.get(game_id);
            if (!game) {
                return;
            }

            var need_password = game.get('game_has_password');
            var password = need_password ? this.$('#js-join-game-password').val() : '';

            var teams_disabled = game.get('teams_disabled') === true;
            var team_id = this.$('#js-join-game-team option:selected').val();

            if (!teams_disabled && !team_id) {
                this.badTeam();
                return;
            }
            this.$(event.currentTarget).addClass('disabled').text('Joining');
            app.Running.Games.joinGame(game_id, password, team_id);
        },
        badPassword: function(){
            this.$('#join_password_block').addClass('has-error');
            this.$('label[for=js-join-game-password]').text('Invalid Password:');
            this.$('.js-join-game-submit').removeClass('disabled').text('Join');
        },
        badTeam: function(){
            this.$('#join_team_block').addClass('has-error');
            this.$('label[for=js-join-game-team]').text('Must Select A Team:');
        },
        fixFields: function(){
            this.$('.has-error').removeClass('has-error');
            this.$('label[for=js-join-game-password]').text('Password:');
            this.$('label[for=js-join-game-team]').text('Team:');
        },
        finishJoin: function(){
            Backbone.history.navigate('my-profile', {
                trigger: true
            });
        },
        noTeams: function(){
            var teamField = this.$('#js-join-game-team');
            teamField.attr('disabled', true);
            teamField.find('#js-team-placeholder').text('No Teams');
        },
        // toggles the password entry field and team entry field on join game
        checkFields: function() {
            // Password field setup
            this.fixFields();
            // get selected game
            var selected = this.$('#js-select-game option:selected');
            var game_id = selected.val();

            var game = this.collection.get(game_id);
            if (!game) {
                return;
            }

            // check if it needs a password
            var need_password = game.get('game_has_password');

            // grab the password field
            var passwordField = this.$('#js-join-game-password');

            // Set it to no password if we don't have one
            var passwordPlaceholder = need_password ? '' : 'No Password';

            // If we do have one mark it as not disabled
            passwordField.attr('disabled', !need_password);
            passwordField.val(passwordPlaceholder);


            if (game.get('teams_disabled') === true) {
                this.noTeams();
                return;
            }

            // Set teams to loading if the game has teams
            var teamField = this.$('#js-join-game-team');
            teamField.find('#js-team-placeholder').text(strings.loading);

            // Set the teams
            var that = this;
            var url = config.WEB_ROOT + 'game/' + game_id + '/team/';
            app.Running.Teams.fetch({
                url:url,
                success:function(teams){
                    // if there are no teams handle that appropriately
                    if (!app.Running.Teams.length){
                        that.noTeams();
                        game.set('teams_disabled', true);
                        return;
                    }
                    // if we have teams mark it as not disabled
                    teamField.attr('disabled', false);

                    // render team select
                    var teamOptionsTemplate = app.Templates["select-game-team-options"];
                    var teamOptionsHTML = teamOptionsTemplate({teams: app.Running.Teams.toJSON()});
                    teamField.html(teamOptionsHTML);
                },
                error:function(err){
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
})();
