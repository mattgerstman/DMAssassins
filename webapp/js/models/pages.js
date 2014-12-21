//
// js/models/pages.js
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
    app.Models.Pages = Backbone.Model.extend({
        defaults: {
            pages: [],
            next: null,
            previous: null
        },
        needPages: function() {
            alert('You must grant access to your pages to use this feature');
            return;
        },
        getPermission: function() {
            var that = this;
            app.Running.FB.login(function(response) {
                var fb_manage_pages = response.authResponse.grantedScopes.search('manage_pages');
                if (fb_manage_pages == -1)
                {
                    return that.needPages();
                }
                app.Running.Permissions.set('manage_pages', true);
                callback();
            }, {
                scope: 'manage_pages',
                auth_type: 'rerequest',
                return_scopes: true
            });
        },
        checkPermission: function(callback) {
            var manage_pages = app.Running.Permissions.get('manage_pages');
            if (manage_pages)
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
            var url = options.url || '/me/accounts';
            if (that === undefined)
            {
                that = this;
            }
            FB.api(url, function(response) {
                console.log(response);
                if (!response || response.error)
                {
                    if (options.error)
                    {
                        return options.error(that, response);
                    }
                }

                that.set('pages', response.data);

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
        getPage: function(i) {
            var pages = this.get('pages');
            if (!pages)
            {
                return false;
            }
            return pages[i];
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
        }
    });
})();
