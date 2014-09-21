//
// js/views/admin-users-view.js
// dmassassins.js
//
// Copyright (c) 2014 Matt Gerstman
// MIT License.
//
// displays user profile


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
    app.Views.AdminUsersView = Backbone.View.extend({


        template: _.template($('#admin-users-template').html()),
        tagName: 'div',

        // The DOM events specific to an item.
        events: {
            'click .ban-user': 'banUserModal',
            'change select.user-team': 'selectChangeTeam',
            'change select.user-role': 'selectChangeRole',
            'click li.team': 'sortByTeam'        
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
            this.listenTo(this.collection, 'sync', this.render);
            this.listenTo(this.collection, 'change', this.makeDraggable)
            this.listenTo(app.Running.Games, 'game-change', this.collection.fetch);
        },
        banUserModal: function(event) {
            var user_name = $(event.currentTarget).data('user-name');
            var user_id = $(event.currentTarget).data('user-id');
            $('#ban_user_submit').data('user-id', user_id);
            $('#ban_user_modal .user-name').text(user_name)
            $('#ban_user_modal').modal();  
        },
        selectChangeTeam: function(event){
            var user_id = $(event.currentTarget).data('user-id');
            var team_id = $(event.currentTarget).find('option:selected').val()
            var team_name = $(event.currentTarget).find('option:selected').text();
            this.addUserToTeam(user_id, team_id, team_name);
            
        },
        selectChangeRole: function(event){
            var user_id = $(event.currentTarget).data('user-id');
            var role_id = $(event.currentTarget).find('option:selected').val()
            return this.changeUserRole(user_id, role_id);
        },
        changeUserRole: function(user_id, role_id){
            // Sorry Taylor, a model for this one is overkill
            var game_id = app.Running.Games.getActiveGameId();
            var url = config.WEB_ROOT + 'game/' + game_id + '/user/' + user_id + '/role/';
            $.post(url, {role: role_id} ,function(){
                $('#role_saved_'+user_id).fadeIn(500, function(){ $(this).fadeOut(2000) })
            });
        },
        addUserToTeam: function(user_id, team_id, team_name, callback) {
            var that = this;
            var team = new app.Models.Team({user_id: user_id, team_id: team_id})
            var user = app.Running.Users.get(user_id);
            return team.save(null, {
                success: function(){
                    that.collection.get(user_id).setProperty('team', team_name);
                    $('#team_saved_'+user_id).fadeIn(500, function(){ $(this).fadeOut(2000) })
                }
            });     
            
            
        },
        makeDraggable: function() {
            var that = this;            
            var startFunc = function(e, ui) {
                ui.helper.find('.user').remove();
                ui.helper.removeClass('user-grid');
                ui.helper.find('.drag-img').removeClass('hide');
                ui.helper.find('.drag-img').animate({
                    width: 50,
                    height: 50                    
                }, 100);
            };
            
            this.$el.find('.user-grid').draggable({
                handle: '.thumbnail',
                connectWith: '#team_list li',
                tolerance: "pointer",
                helper: 'clone',
                forceHelperSize: true,
                zIndex:5000,
                start: startFunc,
                cursorAt: {left:40, top:25}
            })
        },
        makeDroppable: function() {
            var that = this;
            this.$el.find('#team_list li').droppable({
                hoverClass: 'drop-hover',
                tolerance: "pointer",
                drop: function(event, ui) {
                    var user_id = ui.helper.data('user-id');
                    var team_id = $(this).data('team-id');
                    var team_name = $(this).data('team-name');
                    that.addUserToTeam(user_id, team_id, team_name);
                }
            });
        },
        addUser: function(user, extras){
            var userView = new app.Views.AdminUserView(user);
            this.userViews.push(userView);
            this.$el.find('.admin-users-body').append(userView.render(extras).el);
        },
        sortByTeam: function(event) {
            this.team = $(event.currentTarget).data('team-name');
            if ($(event.currentTarget).attr('id') == 'SHOW_ALL')
                this.team = undefined;
                
            if ($(event.currentTarget).attr('id') == 'NO_TEAM')
                this.team = null;
            
            this.render();            
        },
        render: function() {
            while (this.userViews.length)
            {   
                var view = this.userViews.pop();
                view.remove();
            }

            var data = this.collection.models;
            if (this.team !== undefined)
            {
                var that = this;
                data = _.filter(data, function(user){
                    return user.getProperty('team') == that.team;
                });
            }

            var userSort = function(user) {
                return user.getProperty('first_name');
            }

            console.log(data);
            data = _.sortBy(data, userSort);
            
            this.$el.html(this.template());

            var myRole = app.Running.User.getProperty('user_role');
            var that = this;
            var extras = {
                teams: app.Running.Teams.toJSON(),
                roles: AuthUtils.getRolesMapFor(myRole)
            };    
                    
            _.each(data, function(user){
                that.addUser(user, extras);
            })
            this.teams_view.setElement(this.$('#team_list')).render();
            this.makeDraggable();
            this.makeDroppable();
            this.trigger('render');
            return this;
        }
    })
})(jQuery);