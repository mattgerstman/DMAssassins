//
// js/views/nav-game-view.js
// dmassassins.js
//
// Copyright (c) 2014 Matt Gerstman
// MIT License.
//
// handles the game dropdown in the nav


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
    app.Views.AdminUsersTeamsView = Backbone.View.extend({
    
        template: _.template($('#admin-users-teams-template').html()),
        tagName: 'ul',
        initialize: function() {
            this.collection = app.Running.Teams;
            this.listenTo(this.collection, 'fetch', this.render)
            this.listenTo(this.collection, 'change', this.render)
            this.listenTo(this.collection, 'reset', this.render)
            this.listenTo(this.collection, 'add', this.render)
            this.listenTo(this.collection, 'remove', this.render)
        },
        render: function() {
            var teamSort = function(team){
                return team.team_name;
            }
            var data = { teams: _.sortBy(this.collection.toJSON(), teamSort) };
            this.$el.html(this.template(data));
            return this;

        },

    })

})(jQuery);