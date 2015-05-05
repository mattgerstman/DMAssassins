//
// js/views/modal-kill-user-view.js
// dmassassins.js
//
// Copyright (c) 2014 Matt Gerstman
// MIT License.
//
// displays a user in the manage users page

(function() {
    'use strict';
    app.Views.ModalKillUserView = Backbone.View.extend({
        template: app.Templates['modal-kill-user'],
        tagName:'div',
        el: '.js-modal-wrapper',
        events: {
            'click  .js-kill-user-submit'   : 'killUserSubmit'
        },
        initialize: function(model) {
            this.model = model;
        },
        killUserSubmit: function(e) {
            e.preventDefault();
            this.killUser();
            this.hideModal();
        },
        hideModal: function() {
            this.$('.js-modal-kill-user').modal('hide');
        },
        killUser: function() {
            var sendEmail = $('.js-notify-kill-user').is(':checked');
            var data = { send_email: sendEmail };
            this.model.kill(data, null, function(response) {
                if (response.responseText) {
                    alert(response.responseText);
                }
            });

        },
        render: function() {
            var data = this.model.attributes;
            this.$el.html(this.template(data));
            this.$('.js-modal-kill-user').modal();
            return this;
        }
    });
})()
