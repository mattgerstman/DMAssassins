//
// js/views/nav-view.js
// dmassassins.js
//
// Copyright (c) 2014 Matt Gerstman
// MIT License.
//
// handles the nav bar at the top

(function() {
    'use strict';
    app.Views.NavView = Backbone.View.extend({


        template: app.Templates.nav,
        el: '.js-wrapper-nav',

        tagName: 'nav',
        events: {
            'click .js-nav-link': 'select'
        },
        // constructor
        initialize: function() {
            this.NavGameView = new app.Views.NavGameView();
            this.model = new app.Models.Nav();
            this.listenTo(app.Running.User, 'fetch', this.render);
            this.listenTo(app.Running.User, 'change', this.render);
            this.listenTo(app.Running.Games, 'game-change', this.render);
            this.listenTo(this.model, 'change', this.render);
        },

        // if we don't have a target hide that view
        render: function() {

            var data = this.model.attributes;
            data.brand_name = config.BRAND_NAME;
            this.$el.html(this.template(data));
            this.updateHighlight();
            this.NavGameView.setElement(this.$('.js-dropdown-parent-games')).render();
            this.NavGameView.render();
            return this;
        },
        // select an item on the nav bar
        select: function(event) {
            var target = event.currentTarget;
            if ($(target).hasClass('disabled') || $(target).hasClass('dropdown-toggle')) {
                event.preventDefault();
                return;
            }

            this.$('.navbar-collapse.in').collapse('hide');
            this.highlight(target);

        },
        updateHighlight: function() {
            var fragment = Backbone.history.fragment;
            if (fragment === '') {
                fragment = config.DEFAULT_VIEW;
            }
            var selectedElem = this.$el.find('.js-nav-' + fragment);
            this.highlight(selectedElem);

        },
        // highlight an item on the nav bar and unhighlight the rest of them
        highlight: function(elem) {
            if (this.$(elem).hasClass('js-dropdown-parent')) {
                return;
            }

            if (this.$(elem).attr('dropdown')) {
                var dropdown = $(elem).attr('dropdown');
                var parent = '.js-dropdown-parent-'+ dropdown;
                elem = parent;
            }

            this.$('.active').removeClass('active');
            this.$(elem).addClass('active');
            return this;
        }
    });
})();
