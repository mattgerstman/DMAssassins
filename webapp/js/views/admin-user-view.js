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

        template: _.template($('#template-admin-user').html()),
        tagName:'div',
        initialize: function(model){
            this.model = model;
            this.listenTo(this.model, 'change', this.render);
            this.listenTo(this.model, 'fetch', this.render);
            this.listenTo(this.model, 'save', this.render);
        },
        makeDraggable: function(selector) {
            var that = this;
            var startFunc = function(e, ui) {
                ui.helper.find('.js-user').remove();
                ui.helper.removeClass('user-grid');
                ui.helper.find('.js-drag-img').removeClass('hide');
                ui.helper.find('.js-drag-img').animate({
                    width: 50,
                    height: 50
                }, 100);
            };

            if (selector === undefined) {
                selector = '.user-grid';
            }

            this.$el.find(selector).draggable({
                handle: '.js-draggable-photo',
                connectWith: '.js-droppable-team',
                tolerance: "pointer",
                helper: 'clone',
                forceHelperSize: true,
                zIndex:5000,
                start: startFunc,
                cursorAt: {left:40, top:25}
            });
        },
        render: function(extras){
            var data = this.model.attributes;
            for (var key in extras) {
                data[key] = extras[key];
            }
            this.$el.html(this.template(data));

            var game = app.Running.Games.getActiveGame();
            var teams_enabled = false;
            if (game)
            {
                teams_enabled = game.areTeamsEnabled();
            }

            if (teams_enabled)
            {
                var user_id = this.model.get('user_id');
                this.makeDraggable('#js-user-'+user_id);
            }
            return this;
        }
    });
})(jQuery);
