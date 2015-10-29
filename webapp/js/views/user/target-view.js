//
// js/views/target-view.js
// dmassassins.js
//
// Copyright (c) 2014 Matt Gerstman
// MIT License.
//
// target view

(function() {
    'use strict';
    app.Views.TargetView = Backbone.View.extend({


        template: app.Templates.target,
        tagName: 'div',

        // The DOM events specific to an item.
        events: {
            'click .js-target-kill'     : 'kill',
            'click .js-target-picture'  : 'showFullImage',
            'keyup .js-target-secret'   : 'secretKeyup',
            'click .js-get-friends'     : 'getFriends'
        },
        needFriends: function() {
            alert('You must grant access to your friends to use this feature');
            return;
        },
        getFriends:function() {
            var that = this;
            app.Running.FB.login(function(response) {
                if (!response || response.error)
                {
                    return that.needFriends();
                }
                var user_friends = response.authResponse.grantedScopes.search('user_friends');
                if (user_friends === -1)
                {
                    return that.needFriends();
                }
                app.Running.Permissions.set('user_friends', true);
                app.Running.TargetFriendsModel.fetch({reset: true});
            }, {
                scope: 'user_friends',
                auth_type: 'rerequest',
                return_scopes: true
            });
        },
        // loads picture in a modal window
        showFullImage: function() {
            $('.js-modal-target-photo').modal();
        },
        // constructor
        initialize: function() {
            this.model = app.Running.TargetModel;
            this.listenTo(this.model, 'change', this.render);
            this.listenTo(this.model, 'fetch', this.render);
            this.listenTo(this.model, 'set', this.render);
            this.targetFriendsView = new app.Views.TargetFriendsView();

            // if we have no friends find some
            var friends = app.Running.TargetFriendsModel.get('friends');
            if (!friends.length) {
                app.Running.TargetFriendsModel.fetch();
            }

        },
        // kills your target
        kill: function() {
            var secret = this.$el.find('.js-target-secret').val();
            if (!secret) {
                alert("Enter your target's secret to kill them!");
            }
            this.$('.js-target-secret').val('');
            var view = this;
            this.model.destroy({
                headers: {
                    'X-DMAssassins-Secret': secret
                },
                success: function(model, response) {
                    view.model.fetch();
                },
                error: function(model, response){
                    if (status === 401)
                    {
                        alert(response.responseText);
                    }
                }
            });
        },
        secretKeyup: function(e){
            if (e.keyCode === 13) {
                e.preventDefault();
                var secret = this.$('.js-target-secret').val();
                if (!secret) {
                    return;
                }
                this.kill();
            }
        },
        render: function() {
            var data = this.model.attributes;
            data.teams_enabled = app.Running.Games.getActiveGameTeamsEnabled();
            this.$el.html(this.template(data));
            this.targetFriendsView.$el = this.$el.find('.js-target-friends');
            this.targetFriendsView.render();
            return this;
        }
    });
})();
