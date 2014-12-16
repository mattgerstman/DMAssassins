//
// js/models/photos.js
// dmassassins.js
//
// Copyright (c) 2014 Matt Gerstman
// MIT License.
//
// model for target pages

var app = app || {
    Collections: {},
    Models: {},
    Views: {},
    Routers: {},
    Running: {},
    Session: {}
};

(function() {
    'use strict';
    app.Models.Photos = Backbone.Model.extend({
        defaults: {
            photos: [],
        },
        needPhotos: function() {
            alert('You must grant access to your photos to use this feature');
            return;
        },
        getPermission: function() {
            var that = this;
            app.Running.FB.login(function(response) {
                var fb_user_photos = response.authResponse.grantedScopes.search('user_photos');
                if (fb_user_photos == -1)
                {
                    return that.needPhotos();
                }
                app.Running.Permissions.set('user_photos', true);
                callback();
            }, {
                scope: 'user_photos',
                auth_type: 'rerequest',
                return_scopes: true
            });
        },
        checkPermission: function(callback) {
            var user_photos = app.Running.Permissions.get('user_photos');
            if (user_photos)
            {
                return callback();
            }
            return this.getPermission(options, callback, that);
        },
        fetch: function(options) {
            var that = this;
            this.checkPermission(function(){
                that.doFetch(options);
            });
        },
        doFetch: function(options) {
            var that = this;
            options = options || {};
            var url = options.url || '/me/photos/?limit=8';
            if (that === undefined)
            {
                that = this;
            }
            FB.api(url, function(response) {
                if (!response || response.error)
                {
                    if (options.error)
                    {
                        return options.error(that, response);
                    }
                }

                that.set('photos', response.data);

                var next = null;
                var previous = null;
                if (response.paging)
                {
                    next     = response.paging.next     || null;
                    previous = response.paging.previous || null;
                }
                that.set('next', next);
                that.set('previous', previous);

                if (options.success)
                {
                    return options.success(that, response);
                }
                return true;
            });
        },
        next: function(options) {
            var next = this.get('next');
            if (!next) {
                return false;
            }
            options = options || {};
            options.url = next;
            return this.fetch(options);
        },
        previous: function(options) {
            var previous = this.get('previous');
            if (!previous) {
                return false;
            }
            options = options || {};
            options.url = previous;
            return this.fetch(options);
        },
        setPhoto: function(index) {
            var photos = this.get('photos');
            if (!photos[index]) {
                return false;
            }
            var photo = photos[index];
            var images = photos.images;
            var best = null;
            var bestDiff = Math.pow(2, 32);
            _.each(images, function(image) {
                var diff = image.width - 300;
                if (diff < 0) {
                    return;
                }

                if (diff < bestDiff)
                {
                    best = image.source;
                }
            });
        }
    });
})();
