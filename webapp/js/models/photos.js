//
// js/models/photos.js
// dmassassins.js
//
// Copyright (c) 2014 Matt Gerstman
// MIT License.
//
// model for target pages

(function() {
    'use strict';
    app.Models.Photos = Backbone.Model.extend({
        defaults: {
            photos: [],
            next: null,
            previous: null
        },
        needPhotos: function() {
            alert('You must grant access to your photos to use this feature');
            return;
        },
        getPermission: function() {
            var that = this;
            app.Running.FB.login(function(response) {
                var fb_user_photos = response.authResponse.grantedScopes.search('user_photos');
                if (fb_user_photos === -1)
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
            return this.getPermission(callback);
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
            var limit = screen.width > 768 ? 8 : 2;
            var url = options.url || '/me/photos/?limit='+limit;
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
        getPhoto: function(i) {
            var photos = this.get('photos');
            if (!photos)
            {
                return false;
            }
            return photos[i];
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
        getBestPhoto: function(i) {
            var photos = this.get('photos');
            if (!photos[i]) {
                return false;
            }
            var photo = photos[i];
            var images = photo.images;
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
            if (best)
            {
                return best;
            }
            return photo.source;
        },
        setPhoto: function(i) {
            var photo = this.getPhoto(i);
            if (!photo)
            {
                return false;
            }
            var photo_thumb = this.getBestPhoto(i);

            return this.savePhoto(photo.source, photo_thumb);
        },
        setProfilePhoto: function() {
            var facebook_id = app.Running.User.get('facebook_id');
            var photo       = 'https://graph.facebook.com/'+facebook_id+'/picture?width=1000';
            var photo_thumb = 'https://graph.facebook.com/'+facebook_id+'/picture?width=300';
            return this.savePhoto(photo, photo_thumb);
        },
        savePhoto: function(photo, photo_thumb) {
            app.Running.User.setProperty('photo', photo, true);
            app.Running.User.setProperty('photo_thumb', photo_thumb, true);

            app.Running.User.save(null, {
                success:function(model, response) {
                    model.trigger('new_photo');
                }
            });
            return true;

        }
    });
})();
