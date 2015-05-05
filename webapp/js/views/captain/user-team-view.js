//
// js/views/users-teams-view.js
// dmassassins.js
//
// Copyright (c) 2014 Matt Gerstman
// MIT License.
//
// handles the game dropdown in the nav


(function($) {
    'use strict';
    app.Views.AdminUsersTeamView = Backbone.View.extend({

        template: app.Templates["users-team"],
        events: {
            'click  .js-edit-team'          : 'showEditTeamForm',
            'click  .js-cancel-edit-team'   : 'cancelEditTeam',
            'click  .js-cancel-new-team'    : 'cancelNewTeam',
            'click  .js-create-new-team'    : 'createNewTeam',

            'keyup  .js-new-team-name'      : 'newTeamKeypress',


            'click  .js-save-edit-team'     : 'saveEditTeam',
            'click  .js-delete-team'        : 'deleteTeamModal',
            'click  .js-delete-team-submit' : 'deleteTeam',

            'mouseover'   : 'sidebarMouseover',
            'mouseout'    : 'sidebarMouseout'
        },
        initialize: function(model, isAdmin) {
            this.model = model;
            this.listenTo(this.model, 'fetch', this.render);
            this.listenTo(this.model, 'change', this.render);
            this.isAdmin = isAdmin;
        },
        cancelEditTeam: function(e) {
            e.preventDefault();
            this.hideEditTeam(event);
        },
        hideEditTeam: function(event) {
            this.$('.team-display').removeClass('hide');
            this.$('.edit-team-form').addClass('hide');
        },
        deleteTeamModal: function(event) {
            var team = app.Running.Teams.get(team_id);
            var deleteView = new app.Views.ModalDeleteTeamView(team);
            deleteView.render();
        },
        makeDroppable: function() {
            var that = this;
            this.$('li').droppable({
                hoverClass: 'active',
                tolerance: "pointer",
                drop: function(event, ui) {
                    var user_id = ui.helper.data('user-id');
                    var team_id = this.get('team_id');
                    var team_name = this.get('team_name');

                    var user = app.Running.Users.get(user_id);
                    user.changeTeam(
                        team_id,
                        team_name,
                        // success callback
                        function(user, response) {

                        },
                        // error callback
                        function(user, response) {
                            alert(response.responseText);
                        }
                    );
                }
            });
        },
        showEditTeamForm: function(event) {
            event.preventDefault();
            this.$('.edit-team-form').removeClass('hide');
            this.$('.team-display').addClass('hide');

        },
        saveEditTeam: function() {
            var name = this.$('.edit-team-name').val();
            if (name === this.model.get('team_name'))
            {
                this.hideEditTeam();
                return;
            }

            this.model.on('all', function(){
                console.log(arguments);
            })

            this.model.set('team_name', name);
            this.model.save();
        },

        // prevent the entire dom from scrolling when the sidebar is scrolled
        sidebarMouseover: function(e){
            document.body.style.overflow = 'hidden';
        },
        sidebarMouseout: function(){
            document.body.style.overflow = 'auto';
        },
        render: function(extras) {

            var data = this.model.attributes;
            data.is_admin = this.isAdmin;

            this.$el.html(this.template(data));

            this.makeDroppable();
            return this;

        },

    });

})(jQuery);
