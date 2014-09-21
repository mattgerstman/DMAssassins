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
            'click .ban-user': 'banUserModal'
        },
        // constructor
        initialize: function() {
            this.collection = app.Running.Users;
            this.teams_view = new app.Views.AdminUsersTeamsView();
            this.listenTo(this.collection, 'fetch', this.render);
            this.listenTo(this.collection, 'change', this.render);
            this.listenTo(this.collection, 'reset', this.render);
            this.listenTo(this.collection, 'sync', this.render);
            this.listenTo(this.collection, 'add', this.render);
            this.listenTo(app.Running.Games, 'game-change', this.collection.fetch);
        },
        banUserModal: function(event) {
            var user_name = $(event.currentTarget).data('user-name');
            var user_id = $(event.currentTarget).data('user-id');
            $('#ban_user_submit').data('user-id', user_id);
            $('#ban_user_modal .user-name').text(user_name)
            $('#ban_user_modal').modal();  
        },
        addUserToTeam: function(user_id, team_id, team_name) {
            var that = this;
            var team = new app.Models.Team({user_id: user_id, team_id: team_id})
            var user = app.Running.Users.get(user_id);
            team.save(null, {
                success: function(){
                    that.collection.get(user_id).setProperty('team', team_name);
                    that.render()
                }
            });     
            
            
        },
        makeSortable: function() {
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
        render: function() {
            var userSort = function(user) {
                return user.properties.first_name;
            }
            this.$el.html(this.template({users: _.sortBy(this.collection.toJSON(), userSort)}));
            this.teams_view.setElement(this.$('#team_list')).render();
            this.makeSortable();
            return this;
        }
    })
})(jQuery);