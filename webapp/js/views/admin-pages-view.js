//
// js/views/admin-pages-view.js
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
    app.Views.AdminPagesView = Backbone.View.extend({


        template: _.template($('#template-modal-pages').html()),
        tagName: 'div',
        el: '.js-pages-select',
        events: {
            'click .js-select-page'         : 'selectPage',
            'click .js-page-previous'       : 'previousPage',
            'click .js-page-next'           : 'nextPage'
        },
        nextPage: function() {
            this.model.next();
        },
        previousPage: function() {
            this.model.previous();
        },
        selectPage: function(e) {
            var photo = $(e.currentTarget);
            var index = photo.data('index');

            return this.model.setPage(index);
        },
        // constructor
        initialize: function() {
            this.model = new app.Models.Pages();
            this.listenTo(this.model, 'change', this.render);
            this.listenTo(this.model, 'fetch', this.render);
            this.listenTo(this.model, 'destroy', this.destroyCallback);
            this.listenTo(this.model, 'set', this.render);
        },
        render: function() {
            console.log('render yo');
            var data = this.model.attributes;
            this.$el.html(this.template(data));
            return this;
        }
    });
})(jQuery);
