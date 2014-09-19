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
                ui.helper.find('.user_data').hide();
                ui.helper.animate({
                    width: 50,
                    height: 50                    
                }, 100);                
            };
            
            this.$el.find('.sortable').sortable({
                handle: '.thumbnail',
                connectWith: '#team_list li',
                tolerance: "pointer",
                helper: 'clone',
                forceHelperSize: true,
                start: startFunc,
                cursorAt: {left:40, top:25}
            })
            
            this.$el.find('#team_list li').droppable({
                hoverClass: 'drop-hover',
                tolerance: "pointer",
                drop: function(event, ui) {
                    var user_id = ui.helper.data('user_id');
                    var team_id = $(this).data('team_id');
                    var team_name = $(this).data('team_name');
                    that.addUserToTeam(user_id, team_id, team_name);
                }
            });

        },
        render: function() {
            this.$el.html(this.template({users: this.collection.toJSON()} ));
            this.teams_view.setElement(this.$('#team_list')).render();
            this.makeSortable();
            return this;
        }
    })
})(jQuery);