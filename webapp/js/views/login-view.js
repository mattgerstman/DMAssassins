//
// js/views/login-view.js
// dmassassins.js
//
// Copyright (c) 2014 Matt Gerstman
// MIT License.
//
// shows the login screen


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
    app.Views.LoginView = Backbone.View.extend({
    
        template: _.template($('#login-template').html()),

        events: {
            'click .btn-facebook': 'login'
        },

        initialize: function() {
            this.model = app.Session;
        },
        // call the model login function
        login: function() {

            this.model.login();
        },
        // render the login page
        render: function() {
            this.$el.html(this.template(this.model.attributes));
            return this;
        },
    });

})(jQuery);