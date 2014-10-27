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
        login: function(e) {
            var button = $(e.currentTarget);
            var width = button.width();
            button.html('<i class="fa fa-facebook"></i> | Loading...').css('width', width+'px');
            button.attr('disabled', true);
            
            if( navigator.userAgent.match('CriOS') )
            {
                window.open('https://www.facebook.com/dialog/oauth?client_id='+config.APP_ID+'&redirect_uri='+ document.location.href +'&scope=email,user_friends,public_profile', '', null);
                return;
            }
                

            this.model.login();
        },
        // render the login page
        render: function() {
            this.$el.html(this.template(this.model.attributes));
            return this;
        },
    });

})(jQuery);