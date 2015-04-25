//
// js/views/profile-photos-view.js
// dmassassins.js
//
// Copyright (c) 2014 Matt Gerstman
// MIT License.
//
// target view

(function() {
    'use strict';
    app.Views.ProfilePhotosView = Backbone.View.extend({


        template: app.Templates["modal-change-photo"],
        tagName: 'div',
        el: '.js-profile-select',
        events: {
            'click .js-profile-select-photo' : 'selectPhoto',
            'click .js-photo-previous'       : 'previousPhoto',
            'click .js-photo-next'           : 'nextPhoto'
        },
        nextPhoto: function() {
            this.model.next();
        },
        previousPhoto: function() {
            this.model.previous();
        },
        selectPhoto: function(e) {
            var photo = $(e.currentTarget);
            var index = photo.data('index');
            if (index === 'profile')
            {
                return this.setProfilePicture();
            }
            return this.model.setPhoto(index);
        },
        setProfilePicture: function() {
            var wantProfile = confirm('Heads Up! If you pick your profile picture, your assassins photo will always match your current profile picture. If this is what you want click OK.');
            if (!wantProfile)
            {
                return;
            }
            this.model.setProfilePhoto();
        },
        // constructor
        initialize: function() {
            this.model = new app.Models.Photos();
            this.listenTo(this.model, 'change', this.render);
            this.listenTo(this.model, 'fetch', this.render);
            this.listenTo(this.model, 'destroy', this.destroyCallback);
            this.listenTo(this.model, 'set', this.render);
        },
        render: function() {
            var data = this.model.attributes;
            this.$el.html(this.template(data));
            return this;
        }
    });
})();
