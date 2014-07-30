// js/models/nav.js

var app = app || {
  Models: {},
  Views: {},
  Routers: {},
  Running: {},
  Session: {}
};

(function() {
  'use strict';

  app.Models.Nav = Backbone.Model.extend({
    defaults: {
      'left': [
        'Target',
        'My Profile',
        'Leaderboard',
        'Rules'],
      'right': [
        {
          'Admin': [
            'Users',
            'Teams',
            'Plot Twists',
            'Twitter',
            'Game Settings'
          ]
        },
        'Logout'

      ]
    },
    initialize: function() {}
  })
})();