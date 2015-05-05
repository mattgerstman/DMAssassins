//
// js/views/modal-ban-user-view.js
// dmassassins.js
//
// Copyright (c) 2014 Matt Gerstman
// MIT License.
//
// displays a user in the manage users page

(function() {
    'use strict';
    app.Views.ModalBanUserView = Backbone.View.extend({
        template: app.Templates['modal-ban-user'],
        tagName:'div',
        el: '.js-modal-wrapper',
        events: {
            'click  .js-ban-user-submit'    : 'banUserSubmit'
        },
        initialize: function(model) {
            this.model = model;
        },
        banUserSubmit: function(e) {
            e.preventDefault();
            this.hideModal();           
            this.banUser();
        },
        hideModal: function() {
            this.$('.js-modal-ban-user').modal('hide');
        },
        banUser: function() {
            var sendEmail = this.$('.js-notify-ban-user').is(':checked');
            this.model.destroy({
                wait: true,
                headers: {
                    'X-DMAssassins-Send-Email': sendEmail
                },
                url: this.model.url() + 'ban/',
                success: function(){
                    // DROIDS change to a listener and remove the view
                },
                error: function(model, response){
                    alert(response.responseText);
                }
            });

        },
        render: function() {
            var data = this.model.attributes;
            this.$el.html(this.template(data));
            this.$('.js-modal-ban-user').modal();
            return this;
        }
    });
})()
