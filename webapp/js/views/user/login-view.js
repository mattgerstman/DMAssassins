//
// js/views/login-view.js
// dmassassins.js
//
// Copyright (c) 2014 Matt Gerstman
// MIT License.
//
// shows the login screen

(function() {
    'use strict';
    app.Views.LoginView = Backbone.View.extend({
        template: app.Templates.login,
        events: {
            'click .js-login': 'login'
        },
        initialize: function() {
            this.model = app.Session;
            return this;
        },
        // call the model login function
        login: function(e) {
            var button = $(e.currentTarget);
            var width = button.width();
            button.html('<i class="fa fa-facebook"></i> | ' + strings.loading).css('width', width+'px');
            button.attr('disabled', true);
            this.model.login();
            return this;
        },
        // render the login page
        render: function() {
            this.$el.html(this.template(this.model.attributes));
            return this;
        },
    });
})();
