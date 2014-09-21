//
// js/views/admin-user-view.js
// dmassassins.js
//
// Copyright (c) 2014 Matt Gerstman
// MIT License.
//
// displays user profile


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
    app.Views.AdminUserView = Backbone.View.extend({

        template: _.template($('#admin-user-template').html()),
        tagName:'div',
        render: function(){
            this.$el.html(this.template(this.model.attributes))
            return this;
        }    
    })
})(jQuery);
    