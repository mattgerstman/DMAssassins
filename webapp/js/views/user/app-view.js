//
// js/views/app-view.js
// dmassassins.js
//
// Copyright (c) 2014 Matt Gerstman
// MIT License.
//
// loads pages within the body of the app

(function() {
    'use strict';
    app.Views.AppView = Backbone.View.extend({
        el: '#app',
        events: {
            'click .js-support': 'support'
        },
        // constructor
        initialize: function() {
            this.$body = this.$('.js-wrapper-app');
            return this;
        },
        // renders a page within the body of the app
        renderPage: function(page) {
            // Removes modal backdrop if we rapidly change pages
            $('.modal-backdrop').remove();
            this.$body.html(page.render().el);
            return this;
        },
        support: function() {
            var supportView = new app.Views.SupportView();
            supportView.render();
        },
        setCurrentView: function(view) {
            if (app.Running.currentView)
                app.Running.currentView.remove();
            app.Running.currentView = view;
            return this;
        }
    });
})();
