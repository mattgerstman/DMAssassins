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
        initialize: function() {
            this.$body = $('#app_body');
        },
        // renders a page within the body of the app
        renderPage: function(page) {
            this.$body.html(page.render().el);
        },
        setCurrentView: function(view) {
            if (app.Running.currentView)
                app.Running.currentView.remove();
            app.Running.currentView = view;

        }
    })
})(jQuery);