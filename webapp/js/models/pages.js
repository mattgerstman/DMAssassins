//
// js/models/pages.js
// dmassassins.js
//
// Copyright (c) 2014 Matt Gerstman
// MIT License.
//
// model for target pages

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
                if (fb_manage_pages === -1)
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
            var url = options.url || '/me/accounts?limit=9';
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
        },
        setPage: function(i) {
            var page = this.getPage(i);
            if (!page) {
                return false;
            }
            var pageId = page.id;
            var pageAccessToken = page.access_token;
            var pageName = page.name;
            if (!pageId || !pageAccessToken || !pageName) {
                return false;
            }
            var game = app.Running.Games.getActiveGame();
            game.set('game_page_id', pageId);
            game.set('game_page_access_token', pageAccessToken);
            game.set('game_page_name', pageName);
            var url = game.gameUrl();

            if (app.Running.Permissions.get('publish_actions'))
            {
                game.save(null, {
                    url: url,
                    success: function(model, response) {

                    }
                });
                return;
            }

            app.Running.FB.login(function(response) {
                var fb_publish_actions = response.authResponse.grantedScopes.search('publish_actions');
                if (fb_publish_actions === -1)
                {
                    return that.needPages();
                }
                app.Running.Permissions.set('publish_actions', true);
                game.save(null, {
                    url: url,
                    success: function(model, response) {

                    }
                });
            }, {
                scope: 'publish_actions',
                auth_type: 'rerequest',
                return_scopes: true
            });

        }
    });
})();
