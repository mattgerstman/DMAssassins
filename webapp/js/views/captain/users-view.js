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
            'click  .js-new-team-open'      : 'showNewTeam',
            'click  .js-team-name '         : 'sortByTeam',
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
            this.listenTo(app.Running.Teams, 'destroy', this.render);

            this.collection.fetch({reset:true})

        },
        sortByTeam: function(event) {
            event.preventDefault();
            this.team = $(event.currentTarget).data('team-name');
            this.team_id = $(event.currentTarget).data('team-id');
            if (this.team_id === 'SHOW_ALL') {
                this.team = undefined;
                this.team_id = 'all';
            }

            if (this.team_id === 'NO_TEAM') {
                this.team = "null";
                this.team_id = 'null';
            }

            this.render();
        },
        showNewTeam: function(event) {
            event.preventDefault();
            this.$el.find('.js-new-team-open').addClass('hide');
            this.$el.find('.js-form-new-team').removeClass('hide');
            this.$el.find('.js-form-new-team input').focus();
        },
        hideNewTeam: function() {
            this.$el.find('.js-new-team-open').removeClass('hide');
            this.$el.find('.js-form-new-team').addClass('hide');
        },
        cancelNewTeam: function(event) {
            if (event)
                event.preventDefault();
            this.hideNewTeam();
        },
        blurTeamForm: function() {
            var team_name = this.$el.find('.new-team input').val();
            if (!team_name) {
                this.hideNewTeam();
            }

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
            if (this.team !== undefined)
            {
                data = _.filter(data, function(user){
                    return user.getProperty('team') === that.team;
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
