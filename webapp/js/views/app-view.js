//
// js/views/app-view.js
// dmassassins.js
//
// Copyright (c) 2014 Matt Gerstman
// MIT License.
//
// loads pages within the body of the app

var app = app || {
    Collections: {},
    Models: {},
    Views: {},
    Routers: {},
    Running: {},
    Session: {}
};

(function() {
    'use strict';
    app.Views.AppView = Backbone.View.extend({
        el: '#app',
        // constructor
        events: {
            'click .js-support' : 'clickSupport'
        },
        initialize: function() {
            this.$body = this.$('.js-wrapper-app');
            return this;
        },
        // handler for when the user clicks support
        clickSupport: function(e) {
            e.preventDefault();
            this.showSupport();
            return this;
        },
        showSupport: function() {
            var supportView = new app.Views.SupportView();
            supportView.render();
            return this;
        },
        // renders a page within the body of the app
        renderPage: function(page) {
            // Removes modal backdrop if we rapidly change pages
            $('.modal-backdrop').remove();
            this.$body.html(page.render().el);
            return this;
        },
        setCurrentView: function(view) {
            if (app.Running.currentView)
                app.Running.currentView.remove();
            app.Running.currentView = view;
            return this;
        }
    });
})();
