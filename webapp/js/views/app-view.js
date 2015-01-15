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

(function($) {
    'use strict';
    app.Views.AppView = Backbone.View.extend({
        el: '#app',
        // constructor
        events: {
            'click .js-support' : 'support'
        },
        initialize: function() {
            this.$body = $('#js-wrapper-app');
        },
        support: function(e) {
            console.log(e);
            e.preventDefault();
            var supportView = new app.Views.SupportView();
            supportView.render();
        },
        // renders a page within the body of the app
        renderPage: function(page) {
            // Removes modal backdrop if we rapidly change pages
            $('.modal-backdrop').remove();
            this.$body.html(page.render().el);
        },
        setCurrentView: function(view) {
            if (app.Running.currentView)
                app.Running.currentView.remove();
            app.Running.currentView = view;

        }
    });
})(jQuery);
