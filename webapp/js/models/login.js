// js/models/login.js

var app = app || {
  Models: {},
  Views: {},
  Routers: {},
  Running: {},
  Session: {}
};
(function() {
  'use strict';

  app.Models.Login = Backbone.Model.extend({
    defaults: {
      logo: '/assets/img/logo.png'
    },
    login: function() {

      var parent = this;

      FB.getLoginStatus(function(response) {
        if (response.status === 'connected') {
          // Logged into your app and Facebook.

          parent.createSession(response);

          //					var authKey = Base64.encode(userID+':'+token);


        } else if (response.status === 'not_authorized') {
          // The person is logged into Facebook, but not your app.
          FB.login(function(response) {
            parent.createSession(response);
          }, {
            scope: 'public_profile,email,user_friends,user_photos'
          })

        } else {

          FB.login(function(response) {
            parent.createSession(response);
          }, {
            scope: 'public_profile,email,user_friends,user_photos'
          })

          // The person is not logged into Facebook, so we're not sure if
          // they are logged into this app or not.
        }
      })

    },
    createSession: function(response) {
      var url = WEB_ROOT + 'session/';
      var data = {
        'facebook_id': response.authResponse.userID,
        'facebook_token': response.authResponse.accessToken
      }
      console.log(data);
      $.post(url, data, function(response) {
        app.Session.username = response.response.username;
        console.log(response.response.username);
        console.log(app.Session.username);
        app.Running.Router.navigate('target', true)
      });

    }
  })
})();