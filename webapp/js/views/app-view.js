var app = app || {
  Models: {},
  Views: {},
  Routers: {},
  Running: {},
  Session: {}
};

(function($) {
  'use strict';
  app.Views.AppView = Backbone.View.extend({


    el: '#app',

    initialize: function() {
      this.$body = $('#app_body');

    },
    renderPage: function(page) {
      this.$body.html(page.render().el);
    }
  })
})(jQuery);