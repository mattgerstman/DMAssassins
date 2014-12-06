//
// js/views/profile-view.js
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
    app.Views.ProfileView = Backbone.View.extend({


        template: _.template($('#template-profile').html()),
        tagName: 'div',

        // The DOM events specific to an item.
        events: {
            'keyup .js-email'                   : 'emailEnter',
            'click .js-email-settings'          : 'showEmailModal',
            'click .js-email-settings-save'     : 'saveEmailSettings',
            'click .js-profile-picture'         : 'showFullImage',
            'click .js-quit-game'               : 'showQuitModal',
            'click .js-quit-game-confirm'       : 'quitGame',
            'shown.bs.modal.js-email-settings'  : 'focusEmail',
            'shown.bs.modal.js-quit-secret'     : 'focusSecret'
        },
        // load profile picture in modal window
        showFullImage: function() {
            $('.js-modal-profile-photo').modal();
        },
        // load quit confirm modal
        showQuitModal: function() {
            var templateVars = {
                quit_game_name: app.Running.Games.getActiveGame().get('game_name')
            };
            var template = _.template($('#quit-modal-template').html());
            var html = template(templateVars);
            $('.js-wrapper-quit-modal').html(html);
            $('.js-profile-quit-modal').modal();
        },
        focusSecret: function(){
            $('.js-quit-secret').focus();
        },
        quitGame: function() {
            var secret = this.$el.find('#js-quit-secret').val();
            this.model.quit(secret);

        },
        focusEmail: function(){
            $('.js-email').focus();
        },
        showEmailModal: function(){
            $('.js-modal-email-settings').modal();
        },
        emailEnter: function(event) {
            if (event.which == 13) {
                this.saveEmailSettings();
            }
        },
        saveEmailSettings: function(){
            var email = $('.js-email').val();
            var allow_email = $('.js-allow-email').is(':checked') ? 'true' : 'false';
            this.model.saveEmailSettings(email, allow_email);
            $('.js-modal-email-settings').modal('hide');
        },
        // constructor
        initialize: function(params) {
            this.model = app.Running.User;
            this.listenTo(this.model, 'change', this.render);
            this.listenTo(this.model, 'fetch', this.render);
            this.listenTo(this.model, 'destroy', this.destroyCallback);
            this.listenTo(this.model, 'set', this.render);
            this.listenTo(app.Running.Games, 'game-change', this.render);
            this.listenTo(app.Running.Games, 'change', this.render);
            this.listenTo(app.Running.Games, 'join-game', this.render);

        },
        destroyCallback: function() {
            $('.js-profile-quit-modal').modal('hide');
        },
        render: function() {
            $('.modal-backdrop').remove();
            var data = this.model.attributes;
            data.teams_enabled = false;
            var game = app.Running.Games.getActiveGame();
            if (game) {
                data.teams_enabled = game.areTeamsEnabled();
            }
            data.allow_email = data.properties.allow_email == 'true';

            var role = app.Running.User.getProperty('user_role');
            var allow_quit = !AuthUtils.requiresCaptain(role);
            data.allow_quit = allow_quit;
            this.$el.html(this.template(data));
            return this;
        }
    });
})(jQuery);
