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
            'click  .js-edit-team'          : 'clickEditTeam',
            'click  .js-cancel-edit-team'   : 'cancelEditTeam',

            'keyup  .js-edit-team-name'     : 'editTeamKeypress',


            'click  .js-save-edit-team'     : 'clickSaveEdit',
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
        editTeamKeypress: function(e) {
            if (e.keyCode === 27) {
                e.preventDefault();
                this.hideEditTeam();
            }
            if (e.keyCode === 13) {
                e.preventDefault();
                this.saveEditTeam();
            }

        },
        cancelEditTeam: function(e) {
            e.stopPropagation();
            this.hideEditTeam();
        },
        hideEditTeam: function() {
            this.$('.team-display').removeClass('hide');
            this.$('.edit-team-form').addClass('hide');
        },
        deleteTeamModal: function() {
            var deleteView = new app.Views.ModalDeleteTeamView(this.model);
            deleteView.render();
        },
        makeDroppable: function() {
            var that = this;
            this.$('li').droppable({
                hoverClass: 'active',
                tolerance: "pointer",
                drop: function(event, ui) {
                    var user_id = ui.helper.data('user-id');
                    var team_id = that.model.get('team_id');
                    var team_name = that.model.get('team_name');

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
        clickEditTeam: function(e) {
            e.stopPropagation();
            this.showEditTeamForm();
        },
        showEditTeamForm: function() {
            this.$('.edit-team-form').removeClass('hide');
            this.$('.team-display').addClass('hide');

        },
        clickSaveEdit: function(e) {
            e.stopPropagation();
            this.saveEditTeam();
        },
        saveEditTeam: function() {
            var name = this.$('.js-edit-team-name').val();
            if (name === this.model.get('team_name'))
            {
                this.hideEditTeam();
                return;
            }

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
