//
// js/views/admin-plot-twists-view.js
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
    app.Views.AdminPlotTwistsView = Backbone.View.extend({

        template: _.template($('#admin-plot-twists-template').html()),
        tagName:'div',
        events: {
          'click a':'defaultTwistHandler'
        },
        initialize: function(){
            this.model = app.Running.Games.getActiveGame();
            this.listenTo(this.model, 'fetch', this.render);
            this.listenTo(this.model, 'change', this.render);
            this.listenTo(this.model, 'save', this.render);
        },
        defaultTwistHandler: function(e) {
            e.preventDefault();
            console.log(e.currentTarget);
        },
        render: function(){
            var data = this.model.attributes;
            data.teams_enabled = this.model.areTeamsEnabled();
            this.$el.html(this.template(data));
            return this;
        }    
    });
})(jQuery);
    