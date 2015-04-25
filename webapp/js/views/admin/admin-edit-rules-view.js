//
// js/views/admin-edit-rules-view.js
// dmassassins.js
//
// Copyright (c) 2014 Matt Gerstman
// MIT License.
//
// displays rules for a game

(function() {
    'use strict';
    app.Views.AdminEditRulesView = Backbone.View.extend({


        template: app.Templates["edit-rules"],
        tagName: 'div',

        initialize: function(params) {
            this.model = app.Running.RulesModel;

            this.listenTo(this.model, 'change', this.render);
            this.listenTo(this.model, 'fetch', this.render);
            this.listenTo(this.model, 'set', this.render);
        },
        loadEditor: function(){
            var that = this;
            this.$(".js-rules-editor").markdown({
                savable:true,
                saveButtonClass: 'btn btn-md btn-primary',
                footer: '<div class="rules-saved js-saved hide">'+strings.saving+'</div>',
                onSave: function(event) {
                        var rules = event.getContent();
                        that.model.set('rules', rules);
                        that.$('.js-saved').removeClass('hide');
                        that.model.save(null, {success: function(){
                            that.$('.js-saved').text(strings.saved).fadeOut(2000, function(){
                                $(this).text(strings.saving);
                            });

                        }});
                    },
                });
        },
        render: function() {
            var data = this.model.attributes;
            this.$el.html(this.template(data));
            this.loadEditor();
            return this;
        }

    });
})();
