// js/views/user-view.js

var app = app || {
  Models: {},
  Views: {},
  Routers: {},
  Running: {},
  Session: {}
};

(function($) {
  'use strict';
  app.Views.TargetView = Backbone.View.extend({


    template: _.template($('#target-template').html()),

    tagName: 'div',

    // The DOM events specific to an item.
    events: {
      'click .thumbnail': 'showFullImage',
      'click #kill': 'kill'
    },

    showFullImage: function() {
      $('#photoModal').modal()
    },

    initialize: function() {
      var params = {
        'username': app.Session.username,
        'type': 'target'
      }
      this.model = new app.Models.User(params)
      this.model.fetch()
      this.listenTo(this.model, 'change', this.render)
      this.listenTo(this.model, 'fetch', this.render)
      //		  this.listenTo(this.model, 'destroy', )		  		  
    },
    kill: function() {
      var secret = this.$el.find('#secret').val();
      var url = WEB_ROOT + 'users/' + +'/target/'
      var view = this;
      this.model.destroy({
        headers: {
          Secret: secret
        },
        success: function() {
          view.initialize()
          view.render()
        }
      })
      /*
      		var settings = {
      			type: 'POST',
      			data: {secret:secret, _method:'DELETE'},
      			complete:function(response){
      				console.log(response)
      			}
      		}
      		console.log(settings)
      		$.ajax(url,settings)
      */

    },
    render: function() {
      var view = this;
      view.$el.hide()
      view.$el.html(view.template (view.model.attributes));
      view.$el.fadeIn(500);
      return this;
    }
  })

})(jQuery);