//
// js/views/multi-game-view.js
// dmassassins.js
//
// Copyright (c) 2014 Matt Gerstman
// MIT License.
//
// handles game selection

(function() {
    'use strict';
    app.Views.MultiGameView = Backbone.View.extend({
        template: app.Templates["select-game"],
        tagName: 'div',
        events: {
            'click  .js-show-create-game'            : 'showCreateGame',
            'click  .js-show-join-game'              : 'showJoinGame'
        },
        // shows the create game subview
        showCreateGame: function() {
            Backbone.history.navigate('create-game', {
                trigger: true
            });

        },
        // shows the join game subview
        showJoinGame: function() {
            Backbone.history.navigate('join-game', {
                trigger: true
            });
        },
        // cancels the game creation/selection
        goBack: function() {
            if (!!app.Running.Games.getActiveGameId()) {
                app.Running.Router.back();
                return;
            }
        },
        render: function() {
            this.$el.html(this.template());
            return this;
        }
    });
})();
