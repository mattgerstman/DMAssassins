//
// js/views/nav-view.js
// dmassassins.js
//
// Copyright (c) 2014 Matt Gerstman
// MIT License.
//
// handles the nav bar at the top


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
    app.Views.NavView = Backbone.View.extend({


        template: _.template($('#template-nav').html()),
        el: '#js-wrapper-nav',

        tagName: 'nav',

        events: {
            'click .js-nav-link': 'select'
        },

        // constructor
        initialize: function() {
            this.NavGameView = app.Running.NavGameView;
            this.listenTo(app.Running.TargetModel, 'fetch', this.handleTarget);
            this.listenTo(app.Running.TargetModel, 'change', this.handleTarget);
            this.listenTo(app.Running.User, 'fetch', this.render);
            this.listenTo(app.Running.User, 'change', this.render);
            this.listenTo(app.Running.Games, 'game-change', this.render);
        },

        // if we don't have a target hide that view
        render: function() {
            var role = app.Running.User.getProperty('user_role');  
            var data = {};
            data.is_captain = AuthUtils.requiresCaptain(role);
            data.is_admin = AuthUtils.requiresAdmin(role);
            data.is_super_admin = AuthUtils.requiresSuperAdmin(role);
            
            this.$el.html(this.template(data));
            this.handleTarget();
            
            
            var selectedElem = this.$el.find('#js-nav-' + Backbone.history.fragment.replace('_', '-'));
            this.highlight(selectedElem);
            
            if (app.Running.NavGameView)
                app.Running.NavGameView.setElement(this.$('#js-nav-games-dropdown')).render();
            return this;
        },

        // select an item on the nav bar
        select: function(event) {
            var target = event.currentTarget;
            if ($(target).hasClass('disabled') || $(target).hasClass('dropdown-toggle')) {
                event.preventDefault();
                return;
            }
            
            $('.navbar-collapse.in').collapse('hide');
            console.log(target);
            this.highlight(target);

        },

        // highlight an item on the nav bar and unhighlight the rest of them
        highlight: function(elem) {
            if ($(elem).hasClass('js-dropdown-parent')) {
                return;
            }

            if ($(elem).hasClass('dropdown_item')) {
                var dropdown = $(elem).attr('dropdown');
                var parent = '#js-dropdown-parent-'+ dropdown;
                elem = parent;
            }
            
            $('.active').removeClass('active');
            $(elem).addClass('active');
        },
        handleAdmin: function() {
            var role = app.Running.User.getProperty('user_role');  
            var allowed = AuthUtils.requiresCaptain(role);
            if (allowed) {
                $('#js-dropdown-parent-admin').removeClass('hide');
                return;
            }
            $('#js-dropdown-parent-admin').addClass('hide');
            return;
            
        },
        handleTarget: function() {
            var game = app.Running.Games.getActiveGame();
            if (!game)
            {
                this.disableTarget();
                return;
            }
            if (!game.get('game_started'))
            {
                this.disableTarget();
                return;
            }
            
            if (!app.Running.TargetModel.get('user_id'))
            {
                this.disableTarget();
                return;
            }
            
            this.enableTarget();
            return;
        },

        // hides the target nav item
        enableTarget: function() {
            this.$el.find('#js-nav-target').removeClass('disabled');
            this.$el.find('#js-nav-target a').removeClass('disabled');
        },

        // shows the target nav item
        disableTarget: function() {
            this.$el.find('#js-nav-target').addClass('disabled');
            this.$el.find('#js-nav-target a').addClass('disabled');
        }

    });

})(jQuery);