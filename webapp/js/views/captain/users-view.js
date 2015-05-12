//
// js/views/users-view.js
// dmassassins.js
//
// Copyright (c) 2014 Matt Gerstman
// MIT License.
//
// displays user profile


(function($) {
    'use strict';
    app.Views.AdminUsersView = Backbone.View.extend({


        template: app.Templates.users,
        tagName: 'div',

        // The DOM events specific to an item.
        events: {
            'click  .js-team-name': 'clickTeam',
        },
        team: undefined,
        // constructor
        initialize: function() {
            var myRole = app.Running.User.getProperty('user_role');
            this.collection = app.Running.Users;
            this.userViews = [];
            this.teams_view = new app.Views.AdminUsersTeamsView();
            this.listenTo(this.collection, 'fetch', this.render);
            this.listenTo(this.collection, 'reset', this.render);
            this.listenTo(app.Running.Teams, 'destroy', this.renderDestroyTeam);
            this.listenTo(app.Running.Users, 'change', this.userChange);
            this.team_id = 'all';
            this.collection.fetch({reset:true})

        },
        clickTeam: function(e) {
            e.preventDefault();
            var team_id = $(e.currentTarget).data('team-id');
            this.sortByTeam(team_id);
        },
        sortByTeam: function(team_id) {
            console.log(team_id);
            this.team_id = team_id;
            if (this.team_id === 'SHOW_ALL') {
                this.team_id = 'all';
            }

            if (this.team_id === 'NO_TEAM') {
                this.team_id = 'none';
            }

            this.render();
        },

        selectActiveTeam: function() {
            this.$el.find('.active').removeClass('active');
            this.$el.find('#js-nav-team-'+this.team_id).addClass('active');
        },
        addUser: function(user, extras){
            extras.logged_in = false;
            if (user.get('user_id') === app.Session.get('user_id')) {
                extras.logged_in = true;
            }
            var userView = new app.Views.AdminUserView(user);
            this.userViews.push(userView);
            return userView.render(extras).el;
        },
        renderDestroyTeam: function(team) {
            var team_id = team.get('team_id');
            if (this.team_id === team_id) {
                this.team_id = 'all';
            }
            return this.render();
        },
        userChange: function(user) {
            if (user === undefined) {
                return;
            }

            var team_id = user.getProperty('team_id');
            if (this.team_id === 'all') {
                return;
            }

            if (this.team_id === 'none' && team_id) {
                return this.removeUser(user);
            }

            if (team_id !== this.team_id) {
                return this.removeUser(user);
            }
        },

        removeUser: function(user) {
            var user_id = user.get('user_id');
            this.$('#js-user-'+user_id).remove();
        },
        render: function() {
            $('.modal-backdrop').remove();

            var that = this;
            var data = {};
            var game = app.Running.Games.getActiveGame();
            var teams_enabled = false;
            if (game)
            {
                teams_enabled = game.areTeamsEnabled();
            }
            data.teams_enabled = teams_enabled;

            this.$el.html(this.template(data));

            this.teams_view.setElement(this.$('.js-team-list')).render();

            while (this.userViews.length)
            {
                var view = this.userViews.pop();
                view.remove();
            }

            data = this.collection.models;
            if (this.team_id !== 'all')
            {
                data = _.filter(data, function(user) {
                    if (that.team_id === 'none') {
                        return !user.getProperty('team_id');
                    }

                    return user.getProperty('team_id') === that.team_id;
                });
            }
            this.selectActiveTeam();

            var userSort = function(user) {
                return user.getProperty('first_name');
            };

            data = _.sortBy(data, userSort);

            var myRole = app.Running.User.getProperty('user_role');
            var extras = {
                teams: app.Running.Teams.toJSON(),
                roles: AuthUtils.getRolesMapFor(myRole, teams_enabled),
                is_admin: AuthUtils.requiresAdmin(myRole),
                teams_enabled: teams_enabled
            };


            var frag = document.createDocumentFragment();
            var $frag = $(frag);
            _.each(data, function(user) {
                var userEl = that.addUser(user, extras);;
                $frag.append(userEl);
            });

            this.$('.admin-users-body').append($frag);

            return this;
        }
    });
})(jQuery);
