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
    app.Views.AdminUsersTeamsView = Backbone.View.extend({

        template: app.Templates["users-teams"],
        tagName: 'ul',
        events: {
            'click  .js-create-new-team'    : 'createNewTeam',
            'blur   .js-form-new-team input': 'blurTeamForm',

            'keyup  .js-new-team-name'      : 'newTeamKeypress',


            'mouseover'   : 'sidebarMouseover',
            'mouseout'    : 'sidebarMouseout'
        },
        initialize: function() {
            this.collection = app.Running.Teams;
            this.listenTo(this.collection, 'fetch', this.render);
            this.listenTo(this.collection, 'reset', this.render);
            this.listenTo(this.collection, 'add', this.render);
            this.listenTo(this.collection, 'remove', this.render);
        },

        cancelEditTeam: function(e) {
            e.preventDefault();
            this.hideEditTeam(event);
        },
        hideEditTeam: function(event) {
            var team_id = $(event.currentTarget).data('team-id');
            $('#js-nav-team-'+team_id).find('.team-display').removeClass('hide');
            $('#js-nav-team-'+team_id).find('.edit-team-form').addClass('hide');
        },
        newTeamKeypress: function(event) {
            if (event.keyCode === 27) {
                event.preventDefault();
                this.hideNewTeam();
            }
            if (event.keyCode === 13) {
                event.preventDefault();
                this.createNewTeam();
            }

        },
        createNewTeam: function() {
            var team_name = this.$('.new-team input').val();
            if (!team_name) {
                return;
            }

            var game_id = app.Running.Games.getActiveGameId();
            var url = config.WEB_ROOT + 'game/' + game_id + '/team/';
            var that = this;

            app.Running.Teams.create({team_name:team_name}, {
                success: function(response) {
                    that.teams_view.render();
                    that.selectActiveTeam();
                    that.makeDroppable();
                }
            });
        },

        deleteTeamModal: function(event) {
            var team_name = $(event.currentTarget).data('team-name');
            var team_id = $(event.currentTarget).data('team-id');
            var team = app.Running.Teams.get(team_id);
            var deleteView = new app.Views.ModalDeleteTeamView(team);
            deleteView.render();

        },
        deleteTeam: function(event) {
            var team_id = $(event.currentTarget).data('team-id');
            var team_name = $(event.currentTarget).data('team-name');
            var team = app.Running.Teams.get(team_id);
            var that = this;
            team.destroy({success:function(){
                if (that.team === team_name) {
                    that.team = 'null';
                    that.team_id = 'null';
                }
                app.Running.Users.each(function(user){
                    if (user.getProperty('team') === team_name)
                    {
                        user.setProperty('team', 'null');

                    }
                });
                that.render();
            }});
            $('.js-modal-delete-team').modal('hide');
        },
        showEditTeamForm: function(event) {
            event.preventDefault();
            var team_id = $(event.currentTarget).data('team-id');
            $('#js-nav-team-'+team_id).find('.edit-team-form').removeClass('hide');
            $('#js-nav-team-'+team_id).find('.team-display').addClass('hide');

        },

        saveEditTeam: function(event) {
            event.preventDefault();
            var team_id = $(event.currentTarget).data('team-id');
            var team = app.Running.Teams.get(team_id);
            var name = $('#js-nav-team-'+team_id).find('.edit-team-name').val();
            if (name === team.get('team_name'))
            {
                this.hideEditTeam(event);
                return;
            }

            var that = this;
            team.set('team_name', name);
            team.save();
        },

        // prevent the entire dom from scrolling when the sidebar is scrolled
        sidebarMouseover: function(e){
            document.body.style.overflow = 'hidden';
        },
        sidebarMouseout: function(){
            document.body.style.overflow = 'auto';
        },
        render: function() {
            var teamSort = function(team){
                return team.team_name;
            };

            var frag = document.createDocumentFragment();
            var $frag = $(frag);

            var myRole = app.Running.User.getProperty('user_role');
            var isAdmin = AuthUtils.requiresAdmin(myRole);

            var data = {
                teams: this.collection.toJSON(),
                is_admin: isAdmin
            }

            this.$el.html(this.template(data));

            this.collection.each(function(team){
                var view = new app.Views.AdminUsersTeamView(team, isAdmin);
                $frag.append(view.render().el);
            })

            this.$el.append($frag);
            return this;
        }
    });

})(jQuery);
