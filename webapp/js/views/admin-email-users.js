//
// js/views/admin-email-users-view.js
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
    app.Views.AdminEmailUsersView = Backbone.View.extend({


        template: _.template($('#admin-email-users-template').html()),
        tagName: 'div',

        initialize: function(params) {
            this.model = app.Running.UserEmails;
            this.listenTo(this.model, 'change', this.render);
            this.listenTo(this.model, 'fetch', this.render);
            this.listenTo(this.model, 'set', this.render);
        },
        render: function() {
            var data = this.model.attributes;
            this.$el.html(this.template(data));         
            return this;
        }

    });

})(jQuery);