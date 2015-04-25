//
// js/views/support-view.js
// dmassassins.js
//
// Copyright (c) 2014 Matt Gerstman
// MIT License.
//
// displays rules for a game

(function() {
    'use strict';
    app.Views.SupportView = Backbone.View.extend({

        template: app.Templates.support,
        tagName: 'div',
        el: '.js-wrapper-support',
        events: {
            'click .js-support-submit': 'clickSubmit'
        },
        initialize: function() {
            return this;
        },
        clickSubmit: function(e) {
            e.preventDefault();
            this.submit();
            return this;
        },
        submit: function() {

            var name        = this.$('#js-support-name').val();
            var email       = this.$('#js-support-email').val();
            var subject     = this.$('#js-support-subject').val();
            var message     = this.$('#js-support-message').val();

            var data = {
                name:       name,
                email:      email,
                subject:    subject,
                message:    message
            };

            var model = new app.Models.Support();
            model.set(data);
            model.save(null, {
                success: function() {
                    alert('Successfully Submitted!');
                    $('.js-modal-support').modal('hide');
                },
                error: function(a, b, c) {
                    alert('There was an error submitting your issue, please try again later');
                }
            });
            return this;
        },
        render: function() {
            // var data = this.model.attributes;
            this.$el.html(this.template());
            $('.js-modal-support').modal();
            return this;
        }

    });

})();
