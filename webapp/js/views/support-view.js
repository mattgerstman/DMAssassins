//
// js/views/support-view.js
// dmassassins.js
//
// Copyright (c) 2014 Matt Gerstman
// MIT License.
//
// displays rules for a game


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
    app.Views.SupportView = Backbone.View.extend({


        template: _.template($('#template-support').html()),
        tagName: 'div',
        el: '.js-wrapper-support',
        events: {
            'click .js-support-submit': 'submit'
        },
        initialize: function() {

        },
        submit: function(e) {
            e.preventDefault();
            var name        = $('#support-name').val();
            var email       = $('#support-email').val();
            var subject     = $('#support-subject').val();
            var message     = $('#support-message').val();

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
                    console.log(a);
                    console.log(b);
                    console.log(c);
                    alert('There was an error submitting your issue, please try again later');
                }
            });

        },
        render: function() {
            // var data = this.model.attributes;
            this.$el.html(this.template());
            $('.js-modal-support').modal();
            return this;
        }

    });

})(jQuery);
