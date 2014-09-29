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
            var that = this;
            this.$el.find("#rules-editor").markdown({
                savable:true,
                saveButtonClass: 'btn btn-md btn-primary',
                footer: '<div class="saved hide">Saving...</div>',
                onSave: function(event) {
                        var rules = event.getContent();
                        that.model.set('rules', rules);
                        $('.saved').removeClass('hide')
                        that.model.save(null, {success: function(){
                            $('.saved').text('Saved.').fadeOut(2000, function(){
                                $(this).text('Saving...');    
                            });
                            
                        }});
                    },
                })
        },
        render: function() {
            var data = this.model.attributes;
            this.$el.html(this.template(data));         
            this.loadEditor();
            return this;
        }

    })

})(jQuery);