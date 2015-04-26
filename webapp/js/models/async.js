//
// js/models/async.js
// dmassassins.js
//
// Copyright (c) 2014 Matt Gerstman
// MIT License.
//
// Asyncronous code laoder

(function() {
    'use strict';
    app.Models.Async = Backbone.Model.extend({
        defaults: {
            captain: false,
            admin: false,
            superadmin: false
        },
        loadRole: function (key, callback) {
            var that = this;
            var base64Key = app.Session.get('authKey');
            var url = config.WEB_ROOT + 'js/' + key + '/' + config.VERSION + '/DMAssassins-' + key + '.min.js';
            $.ajax({
                url: url,
                cache: true,
                success: function(data, textStatus, jqxhr) {
                    that.set(key, true);
                    if (typeof callback === 'function') {
                        return callback();
                    }
                },
                error: function( jqxhr, settings, exception ) {
                    alert("Error loading "+key+" portal");
                }
            });

        },
        requiresRole: function(key, callback) {
            var hasRole = this.get(key);
            console.log(key+ ': ' + hasRole);
            if (!hasRole) {
                return this.loadRole(key, callback);
            }
            if (typeof callback === 'function') {
                return callback();
            }
            return null;
        },
        requiresSuperAdmin: function(callback) {
            return this.requiresRole('superadmin', callback);
        },
        requiresAdmin: function(callback) {
            return this.requiresRole('admin', callback);
        },
        requiresCaptain: function(callback) {
            return this.requiresRole('captain', callback);
        }
    });
})();
