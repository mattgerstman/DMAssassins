//
// js/views/target-view.js
// dmassassins.js
//
// Copyright (c) 2014 Matt Gerstman
// MIT License.
//
// target view


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
    app.Views.TargetView = Backbone.View.extend({


        template: _.template($('#target-template').html()),
        tagName: 'div',

        // The DOM events specific to an item.
        events: {
            'click .thumbnail': 'showFullImage',
            'click #kill': 'kill'
        },
        // loads picture in a modal window
        showFullImage: function() {
            $('#photoModal').modal();
        },
        // constructor
        initialize: function() {
            this.model = app.Running.TargetModel;
            this.listenTo(this.model, 'change', this.render);
            this.listenTo(this.model, 'fetch', this.render);
            this.listenTo(this.model, 'set', this.render);
        },
        // kills your target
        kill: function() {
            var secret = this.$el.find('#secret').val();
            var view = this;
            this.model.destroy({
                headers: {
                    'X-DMAssassins-Secret': secret
                },
                success: function() {
                    view.model.fetch();
                }
            });
        },
        render: function() {
            var data = this.model.attributes;
            data.teams_enabled = app.Running.Games.getActiveGame().areTeamsEnabled();
            this.$el.html(this.template(data));
            return this;
        }
    });

})(jQuery);