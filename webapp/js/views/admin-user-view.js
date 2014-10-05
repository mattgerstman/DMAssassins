//
// js/views/admin-user-view.js
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
    app.Views.AdminUserView = Backbone.View.extend({

        template: _.template($('#admin-user-template').html()),
        tagName:'div',
        initialize: function(model){
            this.model = model;
            this.listenTo(this.model, 'change', this.render);
            this.listenTo(this.model, 'fetch', this.render);
            this.listenTo(this.model, 'save', this.render);
        },
        render: function(extras){
            var data = this.model.attributes;
            for (var key in extras) {
                data[key] = extras[key];
            }
            this.$el.html(this.template(data));
            return this;
        }    
    });
})(jQuery);
    