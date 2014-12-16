//
// js/views/target-view.js
// dmassassins.js
//
// Copyright (c) 2014 Matt Gerstman
// MIT License.
//
// target view


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
    app.Views.ProfilePhotosView = Backbone.View.extend({


        template: _.template($('#template-modal-change-photo').html()),
        tagName: 'div',
        el: '.js-profile-select',
        // constructor
        initialize: function() {
            this.model = new app.Models.Photos();
            this.listenTo(this.model, 'change', this.render);
            this.listenTo(this.model, 'fetch', this.render);
            this.listenTo(this.model, 'destroy', this.destroyCallback);
            this.listenTo(this.model, 'set', this.render);
        },
        render: function() {
            var data = this.model.attributes;
            this.$el.html(this.template(data));

            return this;
        }
    });
})(jQuery);
