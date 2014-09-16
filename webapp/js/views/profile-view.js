//
// js/views/profile-view.js
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
    app.Views.ProfileView = Backbone.View.extend({


        template: _.template($('#profile-template').html()),
        tagName: 'div',

        // The DOM events specific to an item.
        events: {
            'click .thumbnail': 'showFullImage',
            'click #quit': 'showQuitModal',
            'click #quit_game_confirm': 'quitGame',
        },

        // load profile picture in modal window
        showFullImage: function() {
            $('#photoModal').modal()
        },
        // load quit confirm modal
        showQuitModal: function() {
            var templateVars = {
                quit_game_name: app.Running.Games.getActiveGame().get('game_name')
            }
            var template = _.template($('#quit-modal-template').html())
            var html = template(templateVars);
            $('#quit_modal_wrapper').html(html);
            $('#quit_modal').modal();
        },
        quitGame: function() {
            var secret = this.$el.find('#quit_secret').val();
            this.model.quit(secret);

        },
        // constructor
        initialize: function(params) {
            this.model = app.Running.User;
            this.listenTo(this.model, 'change', this.render)
            this.listenTo(this.model, 'fetch', this.render)
            this.listenTo(this.model, 'destroy', this.destroyCallback)
            this.listenTo(this.model, 'set', this.render)

        },
        destroyCallback: function() {
            $('#quit_modal').modal('hide')
            $('.modal-backdrop').remove();
        },
        render: function() {
            this.$el.html(this.template(this.model.attributes));
            return this;
        }
    })
})(jQuery);