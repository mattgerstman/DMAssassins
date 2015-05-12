//
// js/views/modal-revive-user-view.js
// dmassassins.js
//
// Copyright (c) 2014 Matt Gerstman
// MIT License.
//
// displays a user in the manage users page

(function() {
    'use strict';
    app.Views.ModalReviveUserView = Backbone.View.extend({
        template: app.Templates['modal-revive-user'],
        tagName:'div',
        el: '.js-modal-wrapper',
        events: {
            'click  .js-revive-user-submit'    : 'reviveUserSubmit'
        },
        initialize: function(model) {
            this.model = model;
        },
        reviveUserSubmit: function(e) {
            console.log(e);
            e.preventDefault();
            this.reviveUser();
            this.hideModal();
        },
        hideModal: function() {
            this.$('.js-modal-revive-user').modal('hide');
        },
        reviveUser: function(event) {
            var sendEmail = $('.js-notify-revive-user').is(':checked');
            var data = {
                send_email: sendEmail
            };
            this.model.revive(data, null, function(response) {
                if (response.responseText) {
                    alert(response.responseText);
                }
            });

        },
        render: function() {
            var data = this.model.attributes;
            this.$el.html(this.template(data));
            this.$('.js-modal-revive-user').modal();
            return this;
        }
    });
})()
