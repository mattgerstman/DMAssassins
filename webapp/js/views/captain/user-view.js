//
// js/views/user-view.js
// dmassassins.js
//
// Copyright (c) 2014 Matt Gerstman
// MIT License.
//
// displays a user in the manage users page

(function() {
    'use strict';
    app.Views.AdminUserView = Backbone.View.extend({

        template: app.Templates.user,
        tagName:'div',
        events: {
            'click  .js-user-action'           : 'clickAction',
            'change .js-user-role'             : 'selectChangeRole',
            'change .js-user-team'             : 'selectChangeTeam'
        },
        initialize: function(model) {
            this.model = model;
            this.listenTo(this.model, 'change', this.renderChange);
            this.listenTo(this.model, 'fetch', this.render);
            this.listenTo(this.model, 'save', this.render);
            this.listenTo(this.model, 'destroy', this.remove);
        },
        makeDraggable: function(selector) {
            var that = this;
            var startFunc = function(e, ui) {
                ui.helper.find('.js-user').remove();
                ui.helper.removeClass('user-grid');
                ui.helper.find('.js-drag-img').removeClass('hide');
                ui.helper.find('.js-drag-img').animate({
                    width: 50,
                    height: 50
                }, 100);
            };

            if (selector === undefined) {
                selector = '.user-grid';
            }

            this.$(selector).draggable({
                handle: '.js-draggable-photo',
                connectWith: '.js-droppable-team',
                tolerance: "pointer",
                helper: 'clone',
                forceHelperSize: true,
                zIndex:5000,
                start: startFunc,
                cursorAt: {left:40, top:25}
            });
        },
        clickAction: function(e) {
            e.preventDefault();
            var action = $(e.target).data('action');

            if (action === 'ban') {
                return this.banUserModal();
            }

            if (action === 'kill') {
                return this.killUserModal();
            }

            if (action === 'revive') {
                return this.reviveUserModal();
            }

            return this;

        },
        banUserModal: function() {
            var modal = new app.Views.ModalBanUserView(this.model);
            modal.render();
            return this;
        },
        killUserModal: function() {
            var modal = new app.Views.ModalKillUserView(this.model);
            modal.render();
            return this;
        },
        reviveUserModal: function() {
            var modal = new app.Views.ModalReviveUserView(this.model);
            modal.render();
            return this;

        },
        selectChangeRole: function(event){
            var user_id = $(event.currentTarget).data('user-id');
            var role_id = $(event.currentTarget).find('option:selected').val();
            if (role_id === this.model.getRole()) {
                return;
            }
            return this.changeUserRole(role_id);
        },
        // DROIDs standardize role/team changes
        changeUserRole: function(role_id){
            var that = this;
            this.model.changeRole(role_id, {
                success: function() {
                    that.showSaved('role');
                },
                error: function(user, response) {
                    alert(response.responseText);
                    that.render();
                }
            });
        },
        selectChangeTeam: function(event){
            var user_id = $(event.currentTarget).data('user-id');
            var team_id = $(event.currentTarget).find('option:selected').val();
            var team_name = $(event.currentTarget).find('option:selected').text();
            this.model.changeTeam(
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


        },
        showSaved: function(item) {
            this.$('.js-'+item+'-saved').fadeIn(500, function(){ $(this).fadeOut(2000); });
        },
        renderChange: function(user, key) {
            if (typeof user === 'undefined') {
                return;
            }

            this.render();

            if (typeof user.changed === 'undefined') {
                return;
            }

            if ((key === 'team_id') || (key === 'team')) {
                this.showSaved('team');
            }

            if (key === 'user_role') {
                this.showSaved('role');
            }
            return this;
        },
        render: function(extras) {
            var data = this.model.attributes;
            for (var key in extras) {
                data[key] = extras[key];
            }
            this.$el.html(this.template(data));

            var game = app.Running.Games.getActiveGame();
            var teams_enabled = false;
            if (!game)
            {
                return this;
            }

            teams_enabled = game.areTeamsEnabled();
            if (!teams_enabled)
            {
                return this;
            }

            var user_id = this.model.get('user_id');
            this.makeDraggable('#js-user-'+user_id);

            return this;
        }
    });
})();
