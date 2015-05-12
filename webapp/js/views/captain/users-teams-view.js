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
            'click  .js-new-team-open'      : 'showNewTeam',
            'click  .js-create-new-team'    : 'createNewTeam',
            'click  .js-cancel-new-team'    : 'cancelNewTeam',
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
            this.listenTo(this.collection, 'change', this.sort);
            this.listenTo(this.collection, 'sort', this.render);
        },
        newTeamKeypress: function(e) {
            if (e.keyCode === 27) {
                e.preventDefault();
                this.hideNewTeam();
            }
            if (e.keyCode === 13) {
                e.preventDefault();
                this.createNewTeam();
            }

        },
        showNewTeam: function(event) {
            event.preventDefault();
            this.$el.find('.js-new-team-open').addClass('hide');
            this.$el.find('.js-form-new-team').removeClass('hide');
            this.$el.find('.js-form-new-team input').focus();
        },
        hideNewTeam: function() {
            this.$('.js-new-team-open').removeClass('hide');
            this.$('.js-form-new-team').addClass('hide');
        },
        cancelNewTeam: function() {
            this.hideNewTeam();
        },
        createNewTeam: function() {
            var team_name = this.$('.new-team input').val();
            if (!team_name) {
                return;
            }

            var game_id = app.Running.Games.getActiveGameId();
            var url = config.WEB_ROOT + 'game/' + game_id + '/team/';
            var that = this;

            app.Running.Teams.create({team_name:team_name});
        },
        // prevent the entire dom from scrolling when the sidebar is scrolled
        sidebarMouseover: function(e){
            document.body.style.overflow = 'hidden';
        },
        sidebarMouseout: function(){
            document.body.style.overflow = 'auto';
        },
        sort: function() {
            this.collection.sort();
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
