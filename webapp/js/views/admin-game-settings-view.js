//
// js/views/admin-user-view.js
// dmassassins.js
//
// Copyright (c) 2014 Matt Gerstman
// MIT License.
//
// displays user profile


var app = app || {
    Collections: {},
    Models: {},
    Views: {},
    Routers: {},
    Running: {},
    Session: {}
};

(function($) {
    'use strict';
    app.Views.AdminGameSettingsView = Backbone.View.extend({

        template: _.template($('#admin-game-settings-template').html()),
        tagName:'div',
        events: {
            'click .save-game': 'saveGame',
            'click .start-game': 'startGameModal',
            'click .start-game-submit': 'startGame',
            'click .end-game': 'endGameModal',
            'click .end-game-submit': 'endGame'

        },
        initialize: function(){
            this.model = app.Running.Games.getActiveGame();
            this.listenTo(this.model, 'fetch', this.render)
            this.listenTo(this.model, 'change', this.render)
            this.listenTo(this.model, 'save', this.render)
        },
        saveGame: function(event){
        
            // Get values from form
            var game_name = $('#game_name').val();
            var game_password = $('#game_password').val();
            var game_teams_enabled = $('#teams_enabled').is(':checked') ? 'true' : 'false';
            
            // Set values in model
            this.model.set({
                game_name: game_name,
                game_password: game_password,
                game_teams_enabled: game_teams_enabled},
                {silent:true}
                );    
        
            // Save model
            var url = this.model.gameUrl();
            $(".save-game").text('Saving...');
            this.model.save(null, {
                url: url,
                success: function(model){
                    $(".save-game").text('Saved');        
                    setTimeout(function(){
                        $(".save-game").text('Save');    
                    }, 1000)
                }                
            });
        },
        startGameModal: function(event) {
          $('#start_game_modal').modal();
        },
        startGame: function(event) {
            $('#start_game_modal').modal('hide');
            var that = this;
            var url = this.model.gameUrl();
            $.post(url, function(){
                that.model.set('game_started', true);
            }).error(function(response){
                alert(response.responseText);
            });
        },
        endGameModal: function(event) {
          $('#end_game_modal').modal();
        },
        endGame: function(event) {
            $('#end_game_modal').modal('hide');
            var that = this;
            var url = this.model.gameUrl();

            this.model.destroy({
                url: url,
                success: function() {
                    if (!app.Running.Games.setArbitraryActiveGame()) {
                        Backbone.history.navigate('#logout', {
                            trigger: true
                        });
                        return;
                    }
                }
            });
        },
                
        render: function(){
            $('.modal-backdrop').remove();            
            var data = this.model.attributes;
            data.teams_enabled = data.game_properties.teams_enabled == 'true';
            this.$el.html(this.template(data))
            return this;
        }    
    })
})(jQuery);
    