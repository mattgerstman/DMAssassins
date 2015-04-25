//
// js/views/nav-game-view.js
// dmassassins.js
//
// Copyright (c) 2014 Matt Gerstman
// MIT License.
//
// handles the game dropdown in the nav


(function() {
    'use strict';
    app.Views.NavGameView = Backbone.View.extend({
        template: app.Templates["nav-game"],
        tagName: 'ul',

        events: {
            'click .js-switch-game': 'select'
        },
        // constructor, loads a user id so we can get their games from the model
        initialize: function() {
            this.collection = app.Running.Games;
            this.model = new app.Models.NavGames();
            this.listenTo(this.model, 'change', this.render);

        },
        // loads the items into the dropdown and changes the dropdown title to the current game
        render: function() {
            var data = this.model.attributes;
            this.$el.html(this.template(data));
            return this;

        },
        // select a game from the dropdown
        select: function(event) {
            var game_id = $(event.target).attr('game_id');
            app.Running.Games.setActiveGame(game_id);
        }
    });

})();
