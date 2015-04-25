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
    app.Views.TargetFriendsView = Backbone.View.extend({


        template: app.Templates["target-friends"],
        tagName: 'div',
        el: '.js-target-friends',
        // constructor
        initialize: function() {
            this.model = app.Running.TargetFriendsModel;
            this.listenTo(this.model, 'change', this.render);
            this.listenTo(this.model, 'fetch', this.render);
            this.listenTo(this.model, 'reset', this.render);
            this.listenTo(this.model, 'set', this.render);
            this.listenTo(app.Running.Permissions, 'change', this.render);
        },
        loadMutualFriends: function() {
            var that = this;
            var facebook_id = this.model.get('facebook_id');
            if (!facebook_id) {
                return;
            }
            FB.api('/'+facebook_id+'/friends', {}, function(response) {
                if (response.error) {
                    return;
                }

                if (!response.data.length) {
                    return;
                }

                var friends = [];
                var i = 0;
                var user_facebook_id = app.Running.User.get('facebook_id');

                that.$el.find('.js-mutual-friends').html(this.template({friends: friends}));
            });
        },
        render: function() {
            var data = this.model.attributes;
            data.user_friends = app.Running.Permissions.get('user_friends');
            data.show_button = app.Running.TargetModel.get('game_id');
            this.$el.html(this.template(data));
            return this;
        }
    });
})();
