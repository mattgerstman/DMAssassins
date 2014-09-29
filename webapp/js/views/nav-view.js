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


        template: _.template($('#nav-template').html()),
        el: '#nav_body',

        tagName: 'nav',

        events: {
            'click li a': 'select'
        },

        // constructor
        initialize: function() {
            this.NavGameView = app.Running.NavGameView;
            this.listenTo(app.Running.TargetModel, 'fetch', this.handleTarget)
            this.listenTo(app.Running.TargetModel, 'change', this.handleTarget)
            this.listenTo(app.Running.User, 'fetch', this.render)
            this.listenTo(app.Running.User, 'change', this.render)
            this.listenTo(app.Running.Games, 'game-change', this.render)
        },

        // if we don't have a target hide that view
        render: function() {
            var role = app.Running.User.getProperty('user_role');  
            var data = {};
            data.is_captain = AuthUtils.requiresCaptain(role);
            data.is_admin = AuthUtils.requiresAdmin(role);
            
            this.$el.html(this.template(data));
            this.handleTarget();
            
            var selectedElem = this.$el.find('#nav_' + Backbone.history.fragment);
            this.highlight(selectedElem);
            
            if (app.Running.NavGameView)
                app.Running.NavGameView.setElement(this.$('#games_dropdown')).render();
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
            this.highlight(target)

        },

        // highlight an item on the nav bar and unhighlight the rest of them
        highlight: function(elem) {
            if ($(elem).hasClass('dropdown_parent')) {
                return;
            }

            if ($(elem).hasClass('dropdown_item')) {
                var dropdown = $(elem).attr('dropdown');
                var parent = '#' + dropdown + '_parent';
                elem = parent;
            }
            $('.active').removeClass('active');
            $(elem).addClass('active');
        },
        handleAdmin: function() {
            var role = app.Running.User.getProperty('user_role');  
            var allowed = AuthUtils.requiresCaptain(role);
            if (allowed) {
                $('#admin_parent').removeClass('hide');
                return;
            }
            $('#admin_parent').addClass('hide');
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
            this.$el.find('#nav_target').removeClass('disabled');
            this.$el.find('#nav_target a').removeClass('disabled');
        },

        // shows the target nav item
        disableTarget: function() {
            this.$el.find('#nav_target').addClass('disabled');
            this.$el.find('#nav_target a').addClass('disabled');
        }

    })

})(jQuery);