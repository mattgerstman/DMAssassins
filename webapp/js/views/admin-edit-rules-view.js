//
// js/views/admin-edit-rules-view.js
// dmassassins.js
//
// Copyright (c) 2014 Matt Gerstman
// MIT License.
//
// displays rules for a game


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
    app.Views.AdminEditRulesView = Backbone.View.extend({


        template: _.template($('#admin-edit-rules-template').html()),
        tagName: 'div',

        initialize: function(params) {
            this.model = app.Running.RulesModel;

            this.listenTo(this.model, 'change', this.render)
            this.listenTo(this.model, 'fetch', this.render)
            this.listenTo(this.model, 'set', this.render)
        },
        loadEditor: function(){
            this.$el.find("#rules-editor").markdown()
        },
        render: function() {
            this.$el.html(this.template(this.model.attributes));         
            this.loadEditor();
            return this;
        }

    })

})(jQuery);