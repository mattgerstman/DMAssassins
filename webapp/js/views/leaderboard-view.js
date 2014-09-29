//
// js/views/leaderboard-view.js
// dmassassins.js
//
// Copyright (c) 2014 Matt Gerstman
// MIT License.
//
// shows the list of high scores


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
    app.Views.LeaderboardView = Backbone.View.extend({

        template: _.template($('#leaderboard-template').html()),
        tagName: 'div',

        // constructor
        initialize: function(params) {
            this.model = app.Running.LeaderboardModel;
            this.listenTo(this.model, 'change', this.render)
            this.listenTo(this.model, 'reset', this.render)
            this.listenTo(this.model, 'fetch', this.render)
        },
        // renderer
        render: function() {
            this.$el.html(this.template(this.model.attributes));
            var numCols = 2;
            if (this.model.get('teams_enabled') == 'true')
                numCols = 3;
                
            var options = {
                paging: false,
                searching: false,
                info: false,
                order: [
                    [numCols - 1, 'desc'],
                    [numCols, 'desc']
                ]
            };

            this.$el.find('#user_leaderboard_table').dataTable(options);

            if (this.model.get('teams_enabled')) {
                options.order = [
                    [4, 'desc']
                ]
                this.$el.find('#team_leaderboard_table').dataTable(options);
            }
            return this;
        }

    })

})(jQuery);