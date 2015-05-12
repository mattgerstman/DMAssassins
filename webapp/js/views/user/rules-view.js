//
// js/views/profile-view.js
// dmassassins.js
//
// Copyright (c) 2014 Matt Gerstman
// MIT License.
//
// displays rules for a game


(function() {
    'use strict';
    app.Views.RulesView = Backbone.View.extend({


        template: app.Templates.rules,
        tagName: 'div',

        initialize: function(params) {
            this.model = app.Running.RulesModel;

            this.listenTo(this.model, 'change', this.render);
            this.listenTo(this.model, 'fetch', this.render);
            this.listenTo(this.model, 'set', this.render);
        },

        render: function() {
            var data = this.model.attributes;
            data.rules = marked(data.rules);
            this.$el.html(this.template(data));
            return this;
        }

    });

})();
