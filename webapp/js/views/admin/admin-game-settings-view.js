//
// js/views/admin-game-settings-view.js
// dmassassins.js
//
// Copyright (c) 2014 Matt Gerstman
// MIT License.
//
// displays user profile


(function() {
    'use strict';
    app.Views.AdminGameSettingsView = Backbone.View.extend({

        template: app.Templates["game-settings"],
        tagName:'div',
        events: {
            'click .js-end-game'            : 'endGameModal',
            'click .js-end-game-submit'     : 'endGame',
            'click .js-facebook-page'       : 'facebookPageSetup',
            'click .js-open-plot-twist'     : 'loadTwistModal',
            'click .js-save-game'           : 'saveGame',
            'click .js-start-game'          : 'startGameModal',
            'click .js-start-game-submit'   : 'startGame',
            'click .js-submit-twist'        : 'savePlotTwist',
            'click .js-timer-info'          : 'loadKillTimerModal'
        },
        initialize: function(){
            this.model = app.Running.Games.getActiveGame();
            this.listenTo(this.model, 'fetch', this.render);
            this.listenTo(this.model, 'change', this.render);
            this.listenTo(this.model, 'save', this.render);
        },
        saveGame: function(event){
            event.preventDefault();
            // Get values from form
            var game_name           = $('#js-input-game-name').val();
            var game_password       = $('#js-input-game-password').val();
            var game_teams_enabled  = $('#js-input-teams-enabled').is(':checked') ? 'true' : 'false';
            var game_timezone       = $('#js-timezone').val();

            // Set values in model
            this.model.set({
                    game_name           : game_name,
                    game_password       : game_password,
                    game_teams_enabled  : game_teams_enabled,
                    game_timezone       : game_timezone
                },
                {silent:true}
                );

            // Save model
            var url = this.model.gameUrl();
            $(".js-save-game").text('Saving...');
            this.model.save(null, {
                url: url,
                success: function(model){
                    $(".js-save-game").text('Saved');
                    setTimeout(function(){
                        $(".js-save-game").text('Save');
                    }, 1000);
                }
            });
        },
        startGameModal: function(event) {
            $('.js-modal-start-game').modal();
        },
        startGame: function(event) {
            $('.js-modal-start-game').modal('hide');
            var that = this;
            var sendEmail = $('.js-notify-game-start').is(':checked');
            var data = { send_email: sendEmail };
            this.model.start(data,
                function(model, response) {
                    model.set('game_started', true);
                },
                function(model, response) {
                    if (response.responseText) {
                        alert(response.responseText);
                    }
                }
            );
        },
        endGameModal: function(event) {
            $('.js-modal-end-game').modal();
        },
        endGame: function(event) {
            $('.js-modal-end-game').modal('hide');
            var that = this;
            var url = this.model.gameUrl();
            var sendEmail = $('.js-notify-game-end').is(':checked');

            this.model.destroy({
                url: url,
                headers: {
                    'X-DMAssassins-Send-Email': sendEmail
                },
                success: function(model, response) {
                    alert("The game has successfully ended!\n Thanks for being an admin!");
                    if (!app.Running.Games.setArbitraryActiveGame()) {
                        Backbone.history.navigate('#logout', {
                            trigger: true
                        });
                        return;
                    }
                    Backbone.history.navigate('#my-profile', {
                            trigger: true
                        });
                },
                error: function(model, response) {
                    if (response.responseText) {
                        alert(response.responseText);
                    }
                }
            });
        },
        twistModalOptions: {
            delete_targets: {
                template:     'delete-targets',
                title:        'Delete Targets',
                twist_name:   'delete_targets',
                submit_class: 'btn-primary',
                submit_text:  'Delete Targets',
                checked:       true
            },
            randomize_targets: {
                template:     'randomize-targets',
                title:        'Randomize Targets',
                twist_name:   'randomize_targets',
                submit_class: 'btn-primary',
                submit_text:  'Change Targets',
                checked:       true
            },
            reverse_targets: {
                template:      'reverse-targets',
                title:        'Reverse Targets',
                twist_name:   'reverse_targets',
                submit_class: 'btn-primary',
                submit_text:  'Change Targets',
                checked:       true
            },
            successive_kills: {
                template:      'successive-kills',
                title:        'Kill Mode - Successive Kills Count Double',
                twist_name:   'successive_kills',
                submit_class: 'btn-primary',
                submit_text:  'Enable Plot Twist',
                checked:       false
            },
            strong_weak: {
                template:     'strong-weak',
                title:        'Strong Target Weak',
                twist_name:   'strong_weak',
                submit_class: 'btn-primary',
                submit_text:  'Enable Plot Twist',
                checked:       true
            },
            strong_closed: {
                template:     'strong-closed',
                title:        'Put Strong Players in a Closed Loop',
                twist_name:   'strong_closed',
                submit_class: 'btn-primary',
                submit_text:  'Enable Plot Twist',
                checked:       true
            },
            timer_24: {
                template:     'timer-24',
                title:        '24 Hours To Kill',
                twist_name:   'timer_24',
                submit_class: 'btn-primary',
                submit_text:  'Start Timer',
                checked:       true
            }
        },
        loadTwistModal: function(e){
            e.preventDefault();
            var twist = $(e.currentTarget).attr('id');
            var data = this.twistModalOptions[twist];

            var modal = app.Templates["modal-plot-twist"];

            var detailVars = {};
            detailVars.teams_enabled = this.model.areTeamsEnabled();

            var details = app.Templates.PlotTwist[data.template];
            data.details = details(detailVars);

            var modalHTML = modal(data);
            $('.js-wrapper-modal-plot-twist').html(modalHTML);
            $('.js-modal-plot-twist').modal();

        },
        savePlotTwist: function(e){
            e.preventDefault();

            var button = $(e.currentTarget);
            var data = {};
            data.plot_twist_name  = button.data('twist-name');
            data.send_email       = $('.js-notify-plot-twist').is(':checked');
            var that = this;

            var plotTwist = new app.Models.PlotTwist(data);
            plotTwist.save(null, {
                success: function() {
                    that.model.fetchProperties();
                    alert('The plot twist was launched sucessfully!');
                },
                error: function(model, response) {
                    if (response.responseText === undefined) {
                        alert('There was an error launching the plot twist. Please contact support.');
                        return;
                    }
                    alert(response.responseText);

                }
            });
            $('.js-modal-plot-twist').modal('hide');

        },
        facebookPageSetup: function(e){
            this.pageView = new app.Views.AdminPagesView();
            this.pageView.model.fetch();
            this.pageView.render();
            $('.js-modal-pages').modal();
        },
        loadKillTimerModal: function(e) {
            e.preventDefault();
            $('.js-kill-timer-info').modal();
            var killTimerView = new app.Views.AdminKillTimerView();
            killTimerView.render();

        },
        render: function(){
            $('.modal-backdrop').remove();
            var data = this.model.attributes;
            data.teams_enabled = this.model.getProperty('teams_enabled') === 'true';
            data.has_kill_timer = this.model.getProperty('has_kill_timer') === 'true';

            this.$el.html(this.template(data));

            var timezone = this.model.getProperty('timezone');
            if (timezone) {
                this.$el.find('#js-timezone').val(timezone);
            }
            return this;
        }
    });
})();
