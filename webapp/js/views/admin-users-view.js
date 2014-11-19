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

            'click  .js-ban-user'           : 'banUserModal',
            'click  .js-ban-user-submit'    : 'banUser',
            'click  .js-cancel-edit-team'   : 'cancelEditTeam',            
            'click  .js-cancel-new-team'    : 'cancelNewTeam',
            'click  .js-create-new-team'    : 'createNewTeam',
            'click  .js-delete-team'        : 'deleteTeamModal',
            'click  .js-delete-team-submit' : 'deleteTeam',
            'click  .js-edit-team'          : 'showEditTeamForm',
            'blur   .js-form-new-team input': 'blurTeamForm',
            'click  .js-kill-user'          : 'killUserModal',
            'click  .js-kill-user-submit'   : 'killUser',
            'keyup  .js-new-team-name'      : 'newTeamKeypress',
            'click  .js-new-team-open'      : 'showNewTeam',
            'click  .js-revive-user'        : 'reviveUserModal',
            'click  .js-revive-user-submit' : 'reviveUser',
            'click  .js-save-edit-team'     : 'saveEditTeam',
            'click  .js-team-name '         : 'sortByTeam',
            'change .js-user-team'          : 'selectChangeTeam',
            'change .js-user-role'          : 'selectChangeRole',

            'mouseover .js-sidebar-teams'   : 'sidebarMouseover',
            'mouseout  .js-sidebar-teams'   : 'sidebarMouseout'
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
        },
        banUserModal: function(event) {
            var user_name = $(event.currentTarget).data('user-name');
            var user_id = $(event.currentTarget).data('user-id');
            $('.js-ban-user-submit').data('user-id', user_id);
            $('.js-modal-ban-user .js-user-name').text(user_name);
            $('.js-modal-ban-user').modal();
        },
        killUserModal: function(event) {
            var user_name = $(event.currentTarget).data('user-name');
            var user_id = $(event.currentTarget).data('user-id');
            $('.js-kill-user-submit').data('user-id', user_id);
            $('.js-modal-kill-user .js-user-name').text(user_name);
            $('.js-modal-kill-user').modal();  
        },
        reviveUserModal: function(event) {
            var user_name = $(event.currentTarget).data('user-name');
            var user_id = $(event.currentTarget).data('user-id');
            $('.js-revive-user-submit').data('user-id', user_id);
            $('.js-modal-revive-user .js-user-name').text(user_name);
            $('.js-modal-revive-user').modal();  
        },
          
        banUser: function(event) {
        	var user_id = $(event.currentTarget).data('user-id');
	      	var user = this.collection.get(user_id);
	      	user.destroy({
		      	url: user.url() + 'ban/',
		      	success: function(){
    		        var user_id = user.get('user_id');
    		        $('#js-user-'+user_id).remove();
    		        $('.js-modal-ban-user').modal('hide');      	
		      	},
		      	error: function(response){
                    alert(response.responseText);
		      	}
	      	}); 
		  	
        },
        killUser: function(event) {
        	var user_id = $(event.currentTarget).data('user-id');
	      	var user = this.collection.get(user_id);
            var that = this;
	      	user.kill();
	      	$('.js-modal-kill-user').modal('hide');
        },

        reviveUser: function(event) {
        	var user_id = $(event.currentTarget).data('user-id');
	      	var user = this.collection.get(user_id);
	      	user.revive();
	      	$('.js-modal-revive-user').modal('hide');
        },
        selectChangeTeam: function(event){
            var user_id = $(event.currentTarget).data('user-id');
            var team_id = $(event.currentTarget).find('option:selected').val();
            var team_name = $(event.currentTarget).find('option:selected').text();
            this.addUserToTeam(user_id, team_id, team_name);
            
        },
        selectChangeRole: function(event){
            var user_id = $(event.currentTarget).data('user-id');
            var role_id = $(event.currentTarget).find('option:selected').val();
            return this.changeUserRole(event, user_id, role_id);
        },
        changeUserRole: function(event, user_id, role_id){
            // Sorry Taylor, a model for this one is overkill
            var game_id = app.Running.Games.getActiveGameId();
            var url = config.WEB_ROOT + 'game/' + game_id + '/user/' + user_id + '/role/';            
            $.ajax({
                type:"POST",
                url: url,
                data: {role: role_id},
                success: function(){
                    $('#js-role-saved-'+user_id).fadeIn(500, function(){ $(this).fadeOut(2000); });    
                },
                error: function(response){
                    var originalRole = app.Running.Users.get(user_id).getProperty('user_role');
                    $(event.currentTarget).val(originalRole);
                    alert(response.responseText);
                }
            });
        },
        addUserToTeam: function(user_id, team_id, team_name, callback) {
            var that = this;
            var team = new app.Models.Team({user_id: user_id, team_id: team_id});
            var user = app.Running.Users.get(user_id);
            return team.save(null, {
                success: function(){
                    that.collection.get(user_id).setProperty('team', team_name);
                    $('#js-team-saved-'+user_id).fadeIn(500, function(){ $(this).fadeOut(2000); });
                    if (that.team !== undefined) {
                        $('#js-user-'+user_id).remove();
                    }
                }
            });                         
        },
        makeDroppable: function() {
            var that = this;
            this.$el.find('.js-team-list > li').droppable({
                hoverClass: 'active',
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
                event.preventDefault();
                this.hideNewTeam();                 
             }
             if (event.keyCode == 13) {
                event.preventDefault();
                this.createNewTeam(event); 
             }
             
        },
        showEditTeamForm: function(event) {
            event.preventDefault();
            var team_id = $(event.currentTarget).data('team-id');
            $('#js-nav-team-'+team_id).find('.edit-team-form').removeClass('hide');
            $('#js-nav-team-'+team_id).find('.team-display').addClass('hide');

        },
        hideEditTeam: function(event) {
            var team_id = $(event.currentTarget).data('team-id');
            $('#js-nav-team-'+team_id).find('.team-display').removeClass('hide');
            $('#js-nav-team-'+team_id).find('.edit-team-form').addClass('hide');
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
            var name = $('#js-nav-team-'+team_id).find('.edit-team-name').val();
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
            $('.js-delete-team-submit').data('team-name', team_name);
            $('.js-delete-team-submit').data('team-id', team_id);
            $('.js-delete-team-name').text(team_name);
            $('.js-modal-delete-team').modal();

        },
        deleteTeam: function(event) {
            var team_id = $(event.currentTarget).data('team-id');
            var team_name = $(event.currentTarget).data('team-name');
	      	var team = app.Running.Teams.get(team_id);
	      	var that = this;
	      	team.destroy({success:function(){
	      	    if (that.team == team_name) {
    	      	    that.team = 'null';
    	      	    that.team_id = 'null';
	      	    }
    	      	app.Running.Users.each(function(user){
        	      	if (user.getProperty('team') == team_name)
        	      	{
            	      	user.setProperty('team', 'null');
   
        	      	}
    	      	});
    	      	that.render();    	      	
            }});           
            $('.js-modal-delete-team').modal('hide');
        },
        selectActiveTeam: function() {
            this.$el.find('.active').removeClass('active');
            this.$el.find('#js-nav-team-'+this.team_id).addClass('active');

        },
        sidebarMouseover: function(){
            document.body.style.overflow='hidden';
        },
        sidebarMouseout: function(){
            document.body.style.overflow='auto';
        },
        render: function() {
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
                    return user.getProperty('team') == that.team;
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
                         
            _.each(data, function(user){
                that.addUser(user, extras);
            });

            if (teams_enabled)
            {
//                this.makeDraggable();
                this.makeDroppable();                
            }
            this.trigger('render');
            return this;
        }
    });
})(jQuery);