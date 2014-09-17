//
// js/views/admin-users-view.js
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
    app.Views.AdminUsersView = Backbone.View.extend({


        template: _.template($('#admin-users-template').html()),
        tagName: 'div',

        // The DOM events specific to an item.
        events: {

        },
        // constructor
        initialize: function() {
            this.collection = app.Running.Users;
            this.listenTo(this.collection, 'fetch', this.render);
            this.listenTo(this.collection, 'change', this.render);
            this.listenTo(this.collection, 'reset', this.render);
            this.listenTo(this.collection, 'sync', this.render);
            this.listenTo(this.collection, 'add', this.render);
            this.listenTo(app.Running.Games, 'game-change', this.collection.fetch);
        },
        render: function() {
            this.$el.html(this.template({users: this.collection.toJSON()} ));
            this.$el.find('.sortable').sortable();
            return this;
        }
    })
})(jQuery);