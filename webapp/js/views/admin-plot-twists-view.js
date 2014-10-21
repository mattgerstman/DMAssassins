//
// js/views/admin-plot-twists-view.js
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
    app.Views.AdminPlotTwistsView = Backbone.View.extend({

        template: _.template($('#admin-plot-twists-template').html()),
        tagName:'div',
        events: {
          'click a':'defaultTwistHandler'
        },
        initialize: function(){
            this.model = app.Running.Games.getActiveGame();
            this.listenTo(this.model, 'fetch', this.render);
            this.listenTo(this.model, 'change', this.render);
            this.listenTo(this.model, 'save', this.render);
        },
        defaultTwistHandler: function(e) {
            e.preventDefault();
            var twist = $(e.currentTarget).attr('id');
            switch (twist) {
                case 'targets_randomize':
                    return this.targetsRandomize();
                case 'targets_reverse':
                    return this.targetsReverse();
                case 'targets_strong_weak':
                    return this.targetsStrongWeak();
                case 'targets_strong_closed':
                    return this.targetsStrongClosed();
                case 'kill_mode_normal':
                    return this.killModeNormal();
                case 'kill_mode_successive':
                    return this.killModeSuccessive();
                case 'kill_mode_defend_weak':
                    return this.killModeDefendWeak();
                case 'revive_captains':
                    return this.reviveCaptains();
                case 'revive_strongest':
                    return this.reviveStrongest();
                case 'kill_innocent':
                    return this.killInnocent();
                case 'kill_inactive':
                    return this.killInactive();
            }
        },
        targetsRandomize: function(){
            this.loadTwistModal('#plot-twist-body-targets-randomize-template','Randomize Targets','targets_randomize','btn-primary', 'Change Targets');
        },
        targetsReverse: function(){
            this.loadTwistModal('#plot-twist-body-targets-reverse-template','Reverse Targets','targets_reverse','btn-primary', 'Change Targets');
        },
        targetsStrongWeak: function(){
            this.loadTwistModal('#plot-twist-body-targets-strong-weak-template','Strong Target Weak','targets_strong_weak','btn-primary', 'Change Targets');
        },
        targetsStrongClosed: function(){
            this.loadTwistModal('#plot-twist-body-targets-strong-closed-template','Put Strong Players in a Closed Loop','targets_strong_closed','btn-primary', 'Change Targets');
        },
        killModeNormal: function(){
            this.loadTwistModal('#plot-twist-body-kill-mode-normal-template','Kill Mode - Normal','mill_mode_normal','btn-primary', 'Set Kill Mode');
        },
        killModeSuccessive: function(){
            this.loadTwistModal('#plot-twist-body-kill-mode-successive-kills-template','Kill Mode - Successive Kills Count Double','kill_mode_successive','btn-primary', 'Set Kill Mode');
        },
        killModeDefendWeak: function(){
            this.loadTwistModal('#plot-twist-body-kill-mode-defend-weak-template','Kill Mode - Defend The Weak','kill_mode_successive','btn-primary', 'Set Kill Mode');
        },
        reviveCaptains: function(){
            this.loadTwistModal('#plot-twist-body-revive-team-captains-template','Revive - Team Captains','revive_captains','btn-primary', 'Revive');
        },
        reviveStrongest: function(){
            this.loadTwistModal('#plot-twist-body-revive-strongest-template','Revive - Strongest Players','revive_strongest','btn-primary', 'Revive');
        },
        killInnocent: function(){
            this.loadTwistModal('#plot-twist-body-kill-innocent-template','Kill Players With No Kills','kill_innocent','btn-danger', 'Kill Players');
        },
        killInactive: function(){
            this.loadTwistModal('#plot-twist-body-kill-inactive-template','Kill Players With No Kills in the Past X Hours','kill_inactive','btn-danger', 'Kill Players');
        },
        loadTwistModal: function(templateId, title, submitVal, submitClass, submitText){
            var modal = $('#plot_twist_modal');
            modal.find('.modal-title').text(title);
            modal.find('.twist-submit').val(submitVal).addClass(submitClass).text(submitText);
            var data = {};
            data.teams_enabled = this.model.areTeamsEnabled();
            var details = _.template($(templateId).html());
            modal.find('.twist-details').html(details(data));
            modal.modal();
            
        },
        render: function(){
            var data = this.model.attributes;
            data.teams_enabled = this.model.areTeamsEnabled();
            this.$el.html(this.template(data));
            return this;
        }    
    });
})(jQuery);
    