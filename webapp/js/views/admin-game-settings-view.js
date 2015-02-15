//
// js/views/admin-game-settings-view.js
// dmassassins.js
//
// Copyright (c) 2014 Matt Gerstman
// MIT License.
//
// displays user profile


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
    app.Views.AdminGameSettingsView = Backbone.View.extend({

        template: _.template($('#template-admin-game-settings').html()),
        tagName:'div',
        events: {
            'click .js-end-game'            : 'endGameModal',
            'click .js-end-game-submit'     : 'endGame',
            'click .js-facebook-page'       : 'facebookPageSetup',
            'click .js-open-plot-twist'     : 'loadTwistModal',
            'click .js-save-game'           : 'saveGame',
            'click .js-start-game'          : 'startGameModal',
            'click .js-start-game-submit'   : 'startGame',
            'click .js-submit-twist'        : 'savePlotTwist'
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
            var url = this.model.gameUrl();
            $.post(url, function(){
                that.model.set('game_started', true);
            }).error(function(response){
                alert(response.responseText);
            });
        },
        endGameModal: function(event) {
            $('.js-modal-end-game').modal();
        },
        endGame: function(event) {
            $('.js-modal-end-game').modal('hide');
            var that = this;
            var url = this.model.gameUrl();

            this.model.destroy({
                url: url,
                success: function() {
                    alert("The game has successfully ended!\n Thanks for being an admin!");
                    if (!app.Running.Games.setArbitraryActiveGame()) {
                        Backbone.history.navigate('#logout', {
                            trigger: true
                        });
                        return;
                    }
                    Backbone.history.navigate('#my_profile', {
                            trigger: true
                        });
                }
            });
        },
        twistModalOptions: {
            delete_targets: {
                id:           '#plot-twist-body-delete-targets-template',
                title:        'Delete Targets',
                twist_name:   'delete_targets',
                submit_class: 'btn-primary',
                submit_text:  'Delete Targets',
                checked:       true
            },
            randomize_targets: {
                id:           '#plot-twist-body-randomize-targets-template',
                title:        'Randomize Targets',
                twist_name:   'randomize_targets',
                submit_class: 'btn-primary',
                submit_text:  'Change Targets',
                checked:       true
            },
            reverse_targets: {
                id:           '#plot-twist-body-reverse-targets-template',
                title:        'Reverse Targets',
                twist_name:   'reverse_targets',
                submit_class: 'btn-primary',
                submit_text:  'Change Targets',
                checked:       true
            },
            successive_kills: {
                id:           '#plot-twist-body-successive-kills-template',
                title:        'Kill Mode - Successive Kills Count Double',
                twist_name:   'successive_kills',
                submit_class: 'btn-primary',
                submit_text:  'Enable Plot Twist',
                checked:       false
            },
            strong_weak: {
                id:           '#plot-twist-body-strong-weak-template',
                title:        'Strong Target Weak',
                twist_name:   'strong_weak',
                submit_class: 'btn-primary',
                submit_text:  'Enable Plot Twist',
                checked:       true
            },
            strong_closed: {
                id:           '#plot-twist-body-strong-closed-template',
                title:        'Put Strong Players in a Closed Loop',
                twist_name:   'strong_closed',
                submit_class: 'btn-primary',
                submit_text:  'Enable Plot Twist',
                checked:       true
            },
            timer_24: {
                id:           '#plot-twist-body-timer-24-template',
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

            var modal = _.template($('#template-modal-plot-twist').html());

            var detailVars = {};
            detailVars.teams_enabled = this.model.areTeamsEnabled();

            var details = _.template($(data.id).html());
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
            data.send_email       = $('.js-input-send-twist-email').is(':checked');


            var plotTwist = new app.Models.PlotTwist(data);
            plotTwist.save(null, {
                success: function() {
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
        render: function(){
            $('.modal-backdrop').remove();
            var data = this.model.attributes;
            data.teams_enabled = data.game_properties.teams_enabled == 'true';
            this.$el.html(this.template(data));

            var timezone = this.model.getProperty('timezone');
            if (timezone) {
                this.$el.find('#js-timezone').val(timezone);
            }
            return this;
        }
    });
})(jQuery);
