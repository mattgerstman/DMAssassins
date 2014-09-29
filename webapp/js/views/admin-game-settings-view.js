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
    app.Views.AdminGameSettingsView = Backbone.View.extend({

        template: _.template($('#admin-game-settings-template').html()),
        tagName:'div',
        initialize: function(){
            this.model = app.Running.Games.getActiveGame();
            this.listenTo(this.model, 'change', this.render)
            this.listenTo(this.model, 'fetch', this.render)
            this.listenTo(this.model, 'save', this.render)
        },
        render: function(){
            var data = this.model.attributes;
            this.$el.html(this.template(data))
            return this;
        }    
    })
})(jQuery);
    