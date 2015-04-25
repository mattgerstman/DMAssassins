//
// js/models/nav.js
// dmassassins.js
//
// Copyright (c) 2014 Matt Gerstman
// MIT License.
//
// model for nav

(function() {
    'use strict';
    app.Models.Nav = Backbone.Model.extend({
        defaults: {
            is_captain: false,
            is_admin: false,
            is_super_admin: false,
            target_disabled: false
        },
        initialize: function() {
            this.listenTo(app.Running.User, 'fetch', this.evaluateRole);
            this.listenTo(app.Running.User, 'change', this.evaluateRole);
            this.listenTo(app.Running.Games, 'game-change', this.evaluateRole);
            this.listenTo(app.Running.Games, 'game-change', this.handleTarget);
            this.listenTo(app.Running.TargetModel, 'fetch', this.handleTarget);
            this.listenTo(app.Running.TargetModel, 'change', this.handleTarget);
            this.evaluateRole();
            this.handleTarget();
        },
        evaluateRole: function() {
            var role = app.Running.User.getProperty('user_role');
            var is_captain = AuthUtils.requiresCaptain(role);
            var is_admin = AuthUtils.requiresAdmin(role);
            var is_super_admin = AuthUtils.requiresSuperAdmin(role);
            var data = {
                is_captain: is_captain,
                is_admin: is_admin,
                is_super_admin: is_super_admin
            };
            return this.set(data);
        },
        handleTarget: function() {
            var game = app.Running.Games.getActiveGame();
            if (!game)
            {
                this.set('target_disabled', true);
                return;
            }
            if (!game.get('game_started'))
            {
                this.set('target_disabled', true);
                return;
            }

            if (!app.Running.TargetModel.get('user_id'))
            {
                this.set('target_disabled', true);
                return;
            }

            this.set('target_disabled', false);
            return;
        }
    });
})();
