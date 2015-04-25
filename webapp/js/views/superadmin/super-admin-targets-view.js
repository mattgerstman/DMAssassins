//
// js/views/super-admin-targets-view.js
// dmassassins.js
//
// Copyright (c) 2014 Matt Gerstman
// MIT License.
//
// displays user profile

(function() {
    'use strict';
    app.Views.SuperAdminTargetsView = Backbone.View.extend({
        template: app.Templates.targets,
        tagName:'div',
        initialize: function(){
            this.model = new app.Models.Targets();
            this.listenTo(this.model, 'fetch', this.render);
            this.listenTo(this.model, 'reset', this.render);
            this.listenTo(this.model, 'change', this.render);
        },
        render: function(){
            var data = this.model.attributes;
            this.$el.html(this.template(data));


            var options = {
                paging: false,
                searching: false,
                info: false,
                aaSorting:[],
                aoColumns:[
                    null,
                    null,
                    null,
                    null,
                    null,
                    null,
                    null,
                    null
                ]
            };

            this.$el.find('.js-targets-table').dataTable(options);
            return this;
        }
    });
})();
