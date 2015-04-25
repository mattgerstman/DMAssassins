//
// js/views/profile-view.js
// dmassassins.js
//
// Copyright (c) 2014 Matt Gerstman
// MIT License.
//
// displays user profile

(function() {
    'use strict';
    app.Views.ProfileView = Backbone.View.extend({


        template: app.Templates.profile,
        tagName: 'div',

        // The DOM events specific to an item.
        events: {
            'click .js-change-photo'            : 'changePhotoModal',
            'keyup .js-email'                   : 'emailEnter',
            'click .js-account-settings'          : 'showEmailModal',
            'click .js-account-settings-save'     : 'saveEmailSettings',
            'click .js-profile-picture'         : 'showFullImage',
            'click .js-quit-game'               : 'showQuitModal',
            'click .js-quit-game-confirm'       : 'quitGame',
            'shown.bs.modal.js-account-settings'  : 'focusEmail',
            'shown.bs.modal.js-quit-secret'     : 'focusSecret',
            'hidden.bs.modal'                   : 'render'
        },
        changePhotoModal: function() {
            var that = this;
            this.photosView = new app.Views.ProfilePhotosView();
            this.photosView.model.set('facebook_id', this.model.get('facebook_id'));
            this.photosView.model.fetch();
            this.photosView.render();
            $('.js-modal-profile-change-photo').modal();
        },
        closePhotoModal: function() {
            $('.js-modal-profile-change-photo').modal('hide');
            // this.render();
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
            var template = app.Templates["modal-quit"];
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
            $('.js-modal-account-settings').modal();
        },
        emailEnter: function(event) {
            if (event.which === 13) {
                this.saveEmailSettings();
            }
        },
        saveEmailSettings: function(){
            var email = $('.js-email').val();
            var allow_email = $('.js-allow-email').is(':checked') ? 'true' : 'false';
            var allow_post = $('.js-allow-post').is(':checked') ? 'true' : 'false';
            this.model.set('email', email);
            this.model.setProperty('allow_email', allow_email);
            this.model.setProperty('allow_post', allow_post);
            $('.js-modal-account-settings').modal('hide');
            this.model.save();
        },
        // constructor
        initialize: function(params) {
            this.model = app.Running.User;
            this.listenTo(this.model, 'change', this.render);
            this.listenTo(this.model, 'fetch', this.render);
            this.listenTo(this.model, 'destroy', this.destroyCallback);
            this.listenTo(this.model, 'set', this.render);
            this.listenTo(this.model, 'new_photo', this.closePhotoModal);
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
            data.allow_email = data.properties.allow_email === 'true';

            var role = app.Running.User.getProperty('user_role');
            var allow_quit = !AuthUtils.requiresCaptain(role);
            data.allow_quit = allow_quit;

            data.allow_post = data.properties.allow_post === 'true';
            data.has_page = false;
            if (game)
            {
                if (game.getProperty('game_page_id'))
                {
                    data.has_page = true;
                    data.page_id = game.getProperty('game_page_id');
                    data.page_name = game.getProperty('game_page_name');
                }
            }

            this.$el.html(this.template(data));
            return this;
        }
    });
})();
