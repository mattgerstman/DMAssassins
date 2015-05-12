//
// js/models/permissions.js
// dmassassins.js
//
// Copyright (c) 2014 Matt Gerstman
// MIT License.
//
// User model, manages single user

(function() {
    'use strict';

    app.Models.Permissions = Backbone.Model.extend({
        fetch: function() {
            var that = this;
            FB.api('/me/permissions', function(response){
                var data = {};
                if (!response.data) {
                    return data;
                }
                _.each(response.data, function(permission) {
                    that.set(permission.permission, permission.status === 'granted');
                });
            });
        }
    });
})();
