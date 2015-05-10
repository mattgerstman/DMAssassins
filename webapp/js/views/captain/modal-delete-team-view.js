//
// js/views/modal-delete-team-view.js
// dmassassins.js
//
// Copyright (c) 2014 Matt Gerstman
// MIT License.
//
// displays a user in the manage users page

(function() {
    'use strict';
    app.Views.ModalDeleteTeamView = Backbone.View.extend({
        template: app.Templates['modal-delete-team'],
        tagName:'div',
        el: '.js-modal-wrapper',
        events: {
            'click  .js-delete-team-submit': 'deleteTeamSubmit'
        },
        initialize: function(model) {
            this.model = model;
        },
        deleteTeamSubmit: function(e) {
            e.preventDefault();
            this.hideModal();
            this.deleteTeam();
        },
        hideModal: function() {
            this.$('.js-modal-delete-team').modal('hide');
        },
        deleteTeam: function() {
            var sendEmail = this.$('.js-notify-delete-team').is(':checked');
            var team_name = this.model.get('team_name');
            console.log(team_name);
            this.model.destroy({
                wait: true,
                headers: {
                    'X-DMAssassins-Send-Email': sendEmail
                },
                error: function(model, response){
                    alert(response.responseText);
                }
            });

        },
        render: function() {
            var data = this.model.attributes;
            this.$el.html(this.template(data));
            this.$('.js-modal-delete-team').modal();
            return this;
        }
    });
})()
