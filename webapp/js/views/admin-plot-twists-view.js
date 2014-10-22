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
          'click a':'loadTwistModal',
          'click .twist-submit':'savePlotTwist'
        },
        initialize: function(){
            this.model = app.Running.Games.getActiveGame();
            this.listenTo(this.model, 'fetch', this.render);
            this.listenTo(this.model, 'change', this.render);
            this.listenTo(this.model, 'save', this.render);
        },
        twistModalOptions: {
            targets_randomize: {
                id:           '#plot-twist-body-targets-randomize-template',
                title:        'Randomize Targets',
                twist_name:   'assign_targets',
                twist_value:  'normal',
                submit_class: 'btn-primary',
                submit_text:  'Change Targets',
                checked:       true
            },
            targets_reverse: {
                id:           '#plot-twist-body-targets-reverse-template',
                title:        'Reverse Targets',
                twist_name:   'assign_targets',
                twist_value:  'reverse',
                submit_class: 'btn-primary',
                submit_text:  'Change Targets',
                checked:       true
            },
            targets_strong_weak: {
                id:           '#plot-twist-body-targets-strong-weak-template',
                title:        'Strong Target Weak',
                twist_name:   'assign_targets',
                twist_value:  'strong_weak',
                submit_class: 'btn-primary',
                submit_text:  'Change Targets',
                checked:       true
            },
            targets_strong_closed: {
                id:           '#plot-twist-body-targets-strong-closed-template',
                title:        'Put Strong Players in a Closed Loop',
                twist_name:   'assign_targets',
                twist_value:  'closed_strong',
                submit_class: 'btn-primary',
                submit_text:  'Change Targets',
                checked:       true
            },
            kill_mode_normal: {
                id:           '#plot-twist-body-kill-mode-normal-template',
                title:        'Kill Mode - Normal',
                twist_name:   'kill_mode',
                twist_value:  'normal',
                submit_class: 'btn-primary',
                submit_text:  'Set Kill Mode',
                checked:       false
            },
            kill_mode_successive: {
                id:           '#plot-twist-body-kill-mode-successive-kills-template',
                title:        'Kill Mode - Successive Kills Count Double',
                twist_name:   'kill_mode',
                twist_value:  'successive_kills',
                submit_class: 'btn-primary',
                submit_text:  'Set Kill Mode',
                checked:       false
            },
            kill_mode_defend_weak: {
                id:           '#plot-twist-body-kill-mode-defend-weak-template',
                title:        'Kill Mode - Defend The Weak',
                twist_name:   'kill_mode',
                twist_value:  'defend_weak',
                submit_class: 'btn-primary',
                submit_text:  'Set Kill Mode',
                checked:       false
            },
            revive_captains: {
                id:           '#plot-twist-body-revive-team-captains-template',
                title:        'Revive - Team Captains',
                twist_name:   'revive_strongest',
                twist_value:  '',
                submit_class: 'btn-primary',
                submit_text:  'Revive',
                checked:       false
            },
            revive_strongest: {
                id:           '#plot-twist-body-revive-strongest-template',
                title:        'Revive - Strongest Players',
                twist_name:   'revive_captains',
                twist_value:  '',
                submit_class: 'btn-primary',
                submit_text:  'Revive',
                checked:       false
            },
            kill_innocent: {
                id:           '#plot-twist-body-kill-innocent-template',
                title:        'Kill Players With No Kills',
                twist_name:   'kill_innocent',
                twist_value:  '',
                submit_class: 'btn-danger',
                submit_text:  'Kill Players',
                checked:       true
            },
            kill_inactive: {
                id:           '#plot-twist-body-kill-inactive-template',
                title:        'Kill Players With No Kills in the Past X Hours',
                twist_name:   'kill_inactive',
                twist_value:  '',
                submit_class: 'btn-danger',
                submit_text:  'Kill Players',
                checked:       true
            }
        },
        loadTwistModal: function(e){
            e.preventDefault();
            var twist = $(e.currentTarget).attr('id');            
            var data = this.twistModalOptions[twist];
            
            var modal = _.template($('#plot-twist-modal-template').html());
            
            var detailVars = {};
            detailVars.teams_enabled = this.model.areTeamsEnabled();

            var details = _.template($(data.id).html());
            data.details = details(detailVars);
            
            var modalHTML = modal(data);
            $('#plot-twist-modal-container').html(modalHTML);
            $('#plot-twist-modal').modal();
            
        },
        savePlotTwist: function(e){
            e.preventDefault();
            
            var button = $(e.currentTarget);
            var data = {};
            data.plot_twist_name  = button.data('twist-name');
            data.plot_twist_value = button.data('twist-value');
            data.send_email       = $('#send-twist-email').is(':checked');
            
            var override = $('#plot-twist-value-override').val();
            if (!!override) {
                data.plot_twist_value = override;    
            }
            var plotTwist = new app.Models.PlotTwist(data);
            plotTwist.save();    
            $('#plot-twist-modal').modal('hide');
            
        },
        render: function(){
            var data = this.model.attributes;
            data.teams_enabled = this.model.areTeamsEnabled();
            this.$el.html(this.template(data));
            return this;
        }    
    });
})(jQuery);
    