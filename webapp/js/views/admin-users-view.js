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
            'click .kill-user': 'killUserModal',
            'click .revive-user': 'reviveUserModal',
            'click .ban-user-submit': 'banUser',
            'click .kill-user-submit': 'killUser',
            'click .revive-user-submit': 'reviveUser',
            'change select.user-team': 'selectChangeTeam',
            'change select.user-role': 'selectChangeRole',
            'click  a.team-name ': 'sortByTeam',
            'click .new-team-open': 'showNewTeam',
            'click .create-new-team': 'createNewTeam',
            'click .cancel-new-team': 'cancelNewTeam',
            'keyup .new-team-name': 'newTeamKeypress',
            'blur .new-team-form input': 'blurTeamForm',
            'click .edit-team': 'showEditTeamForm',
            'click .cancel-edit-team': 'cancelEditTeam',
            'click .save-edit-team': 'saveEditTeam',
            'click .delete-team': 'deleteTeamModal',
            'click .delete-team-submit': 'deleteTeam'
        },
        team: undefined,
        // constructor
        initialize: function() {
            var myRole = app.Running.User.getProperty('user_role');
            this.collection = app.Running.Users;
            this.userViews = [];
            this.teams_view = new app.Views.AdminUsersTeamsView();
            this.listenTo(this.collection, 'fetch', this.render);
            this.listenTo(this.collection, 'sync', this.render);
            this.listenTo(this.collection, 'change', this.render)
			this.listenTo(this.collection, 'remove', this.render)
        },
        banUserModal: function(event) {
            var user_name = $(event.currentTarget).data('user-name');
            var user_id = $(event.currentTarget).data('user-id');
            $('.ban-user-submit').data('user-id', user_id);
            $('#ban_user_modal .user-name').text(user_name)
            $('#ban_user_modal').modal();
        },
        killUserModal: function(event) {
            var user_name = $(event.currentTarget).data('user-name');
            var user_id = $(event.currentTarget).data('user-id');
            $('.kill-user-submit').data('user-id', user_id);
            $('#kill_user_modal .user-name').text(user_name)
            $('#kill_user_modal').modal();  
        },
        reviveUserModal: function(event) {
            var user_name = $(event.currentTarget).data('user-name');
            var user_id = $(event.currentTarget).data('user-id');
            $('.revive-user-submit').data('user-id', user_id);
            $('#revive_user_modal .user-name').text(user_name)
            $('#revive_user_modal').modal();  
        },
               
        banUser: function(event) {
        	var user_id = $(event.currentTarget).data('user-id');
	      	var user = this.collection.get(user_id);
	      	user.ban()
		  	$('#ban_user_modal').modal('hide')
        },
        killUser: function(event) {
        	var user_id = $(event.currentTarget).data('user-id');
	      	var user = this.collection.get(user_id);
	      	user.kill()
	      	$('#kill_user_modal').modal('hide')
        },

        reviveUser: function(event) {
        	var user_id = $(event.currentTarget).data('user-id');
	      	var user = this.collection.get(user_id);
	      	user.revive()
	      	$('#revive_user_modal').modal('hide')
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
                    if (that.team !== undefined) {
                        $('#user_'+user_id).remove();
                    }
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
            this.$el.find('#team_list li.team-droppable').droppable({
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
            extras.logged_in = false;
            if (user.get('user_id') == app.Session.get('user_id')) {
                extras.logged_in = true;
            }
            var userView = new app.Views.AdminUserView(user);
            this.userViews.push(userView);
            this.$el.find('.admin-users-body').append(userView.render(extras).el);
        },
        sortByTeam: function(event) {
            event.preventDefault();        
            this.team = $(event.currentTarget).data('team-name');
            this.team_id = $(event.currentTarget).data('team-id');
            if (this.team_id == 'SHOW_ALL') {
                this.team = undefined;
                this.team_id = 'all';
            }
                
                
            if (this.team_id == 'NO_TEAM') {
                this.team = "null";
                this.team_id = 'null'
            }
                
            
            this.render();            
        },
        showNewTeam: function(event) {            
            event.preventDefault();
            this.$el.find('.new-team-open').addClass('hide');
            this.$el.find('.new-team-form').removeClass('hide');
            this.$el.find('.new-team-form input').focus();
        },
        hideNewTeam: function() {             
            this.$el.find('.new-team-open').removeClass('hide');
            this.$el.find('.new-team-form').addClass('hide');
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
        createNewTeam: function(event) {
            if (event)           
                event.preventDefault();
                
            var team_name = this.$el.find('.new-team input').val();
            if (!team_name) {                        
                return;
            }
            var game_id = app.Running.Games.getActiveGameId();
            var url = config.WEB_ROOT + 'game/' + game_id + '/team/';
            var that = this;
            $.post(url, {team_name:team_name}, function(team){
                app.Running.Teams.add(team);
                that.teams_view.render();
                that.selectActiveTeam();
                that.makeDroppable();
            });

        },
        newTeamKeypress: function(event) {
             if (event.keyCode == 27) {
                 this.hideNewTeam();                 
             }
             if (event.keyCode == 13) {
                 this.createNewTeam(event); 
             }
             
        },
        showEditTeamForm: function(event) {
            event.preventDefault();
            var team_id = $(event.currentTarget).data('team-id');
            $('#nav_team_'+team_id).find('.edit-team-form').removeClass('hide');
            $('#nav_team_'+team_id).find('.team-display').addClass('hide');

        },
        hideEditTeam: function(event) {
            var team_id = $(event.currentTarget).data('team-id');
            $('#nav_team_'+team_id).find('.team-display').removeClass('hide');
            $('#nav_team_'+team_id).find('.edit-team-form').addClass('hide');
        },
        cancelEditTeam: function(event) {
            if (event)           
                event.preventDefault();
            this.hideEditTeam(event);
        },
        saveEditTeam: function(event) {
            event.preventDefault();
            var team_id = $(event.currentTarget).data('team-id');
            var team = app.Running.Teams.get(team_id);
            var name = $('#nav_team_'+team_id).find('.edit-team-name').val()
            if (name == team.get('team_name'))
            {
                this.hideEditTeam(event);
                return;
            }
                
            var that = this;
            team.set('team_name', name);
            team.save();
        },
        deleteTeamModal: function(event) {
            event.preventDefault();
            var team_name = $(event.currentTarget).data('team-name');
            var team_id = $(event.currentTarget).data('team-id');
            $('.delete-team-submit').data('team-name', team_name);
            $('.delete-team-submit').data('team-id', team_id);
            $('#delete_team_modal .team-name').text(team_name)
            $('#delete_team_modal').modal();

        },
        deleteTeam: function(event) {
            var team_id = $(event.currentTarget).data('team-id');
            var team_name = $(event.currentTarget).data('team-name');
	      	var team = app.Running.Teams.get(team_id);
	      	var that = this;
	      	team.destroy({success:function(){
	      	    if (that.team == team_name) {
    	      	    that.team = 'null'
    	      	    that.team_id = 'null'
	      	    }
    	      	app.Running.Users.each(function(user){
        	      	if (user.getProperty('team') == team_name)
        	      	{
            	      	user.setProperty('team', 'null'), true;
   
        	      	}
    	      	})
    	      	that.render();    	      	
            }});           
            $('#delete_team_modal').modal('hide');
        },
        selectActiveTeam: function() {
            this.$el.find('.active').removeClass('active');
            this.$el.find('#nav_team_'+this.team_id).addClass('active');

        },
        render: function() {			
			var data = {};
			var game = app.Running.Games.getActiveGame();
			var teams_enabled = false
			if (game)
			{
    			teams_enabled = game.areTeamsEnabled();
			}
            data.teams_enabled = teams_enabled;
			
            this.$el.html(this.template(data));
            
            this.teams_view.setElement(this.$('#team_list')).render();
                    
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
            this.selectActiveTeam();

            var userSort = function(user) {
                return user.getProperty('first_name');
            }

            data = _.sortBy(data, userSort);        

            var myRole = app.Running.User.getProperty('user_role');
            var that = this;
            var extras = {
                teams: app.Running.Teams.toJSON(),
                roles: AuthUtils.getRolesMapFor(myRole, teams_enabled),
                is_admin: AuthUtils.requiresAdmin(myRole),
                teams_enabled: teams_enabled
            };    
                         
            _.each(data, function(user){
                that.addUser(user, extras);
            })

            if (teams_enabled)
            {
                this.makeDraggable();
                this.makeDroppable();                
            }
            this.trigger('render');
            return this;
        }
    })
})(jQuery);