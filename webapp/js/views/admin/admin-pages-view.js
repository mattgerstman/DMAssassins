//
// js/views/admin-pages-view.js
// dmassassins.js
//
// Copyright (c) 2014 Matt Gerstman
// MIT License.
//
// target view


(function() {
    'use strict';
    app.Views.AdminPagesView = Backbone.View.extend({


        template: app.Templates['modal-pages'],
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
            var data = this.model.attributes;
            this.$el.html(this.template(data));
            return this;
        }
    });
})();
