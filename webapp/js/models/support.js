//
// js/models/target.js
// dmassassins.js
//
// Copyright (c) 2014 Matt Gerstman
// MIT License.
//
// model for target pages

(function() {
    'use strict';
    app.Models.Support = Backbone.Model.extend({
        defaults: {
            'name': '',
            'email': '',
            'subject': '',
            'message': ''
        },
        url: function() {
            return config.WEB_ROOT + 'support/' ;
        }
    });
})();
