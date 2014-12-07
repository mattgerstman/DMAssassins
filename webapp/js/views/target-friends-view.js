//
// js/views/target-view.js
// dmassassins.js
//
// Copyright (c) 2014 Matt Gerstman
// MIT License.
//
// target view


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
    app.Views.TargetFriendsView = Backbone.View.extend({


        template: _.template($('#template-target-friends').html()),
        tagName: 'div',
        el: '.js-target-friends',
        // constructor
        initialize: function() {
            this.model = app.Running.TargetFriendsModel;
            this.listenTo(this.model, 'change', this.render);
            this.listenTo(this.model, 'fetch', this.render);
            this.listenTo(this.model, 'set', this.render);
        },
        loadMutualFriends: function() {
            var that = this;
            var facebook_id = this.model.get('facebook_id');
            if (!facebook_id) {
                return;
            }
            FB.api('/'+facebook_id+'/friends', {}, function(response) {
                console.log(response);
                if (response.error) {
                    return;
                }

                if (!response.data.length) {
                    return;
                }

                var friends = [];
                var i = 0;
                var user_facebook_id = app.Running.User.get('facebook_id');

                var template = _.template($('#template-mutual-friends').html());
                that.$el.find('.js-mutual-friends').html(template({friends: friends}));
            });
        },
        render: function() {
            var data             = {};
            var friends          = [];
            var modelData        = this.model.attributes;
            var user_facebook_id = app.Running.User.get('facebook_id');
            var i = 0;
            _.each(modelData.friends, function(friend) {
                if (i > 5) {
                    return;
                }

                if (friend.id != user_facebook_id) {
                    friends.push(friend);
                    i++;
                }
            });
            data.friends = friends;
            data.count   = modelData.count;
            this.$el.html(this.template(data));

            return this;
        }
    });
})(jQuery);
