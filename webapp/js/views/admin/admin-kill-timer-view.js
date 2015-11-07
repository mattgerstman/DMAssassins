//
// js/views/kill-timer-view.js
// dmassassins.js
//
// Copyright (c) 2014 Matt Gerstman
// MIT License.
//
// target view


(function() {
    'use strict';
    app.Views.AdminKillTimerView = Backbone.View.extend({


        template: app.Templates["modal-kill-timer"],
        tagName: 'div',
        el: '.js-wrapper-kill-timer',
        events: {
            'click .js-cancel-timer' : 'cancelTimer'
        },
        // constructor
        initialize: function() {
            var game_id = app.Running.Games.getActiveGameId();
            this.model = new app.Models.KillTimer({game_id: game_id});
            this.model.fetch();
            this.listenTo(this.model, 'change', this.render);
            this.listenTo(this.model, 'fetch', this.render);
            this.listenTo(this.model, 'reset', this.render);
            this.listenTo(this.model, 'set', this.render);
        },
        cancelTimer: function(e) {
            e.preventDefault();
            var check = confirm("Are you sure you want to cancel this kill timer?");
            if (!check) {
                return;
            }
            var sendEmail = $('js-notify-cancel-timer').is(':checked');
            $('.js-kill-timer-info').modal('hide');
            this.model.destroy({
                headers: {
                    'X-DMAssassins-Send-Email': sendEmail
                },
                success: function(model, response) {
                    app.Running.Games.getActiveGame().fetchProperties();
                    alert("Kill timer disabled!");
                },
                error: function(model, response) {
                    if (response.responseText === undefined) {
                        alert('There was an error canceling the kill timer. Please contact support.');
                        return;
                    }
                    alert(response.responseText);
                }
            });
        },
        render: function() {
            var data = this.model.attributes;

            var timer_execute_ts  = this.model.get('execute_ts');
            var timer_min_kill_ts = this.model.get('create_ts');

            var execute_ts  = new Date(timer_execute_ts  * 1000);
            var min_kill_ts = new Date(timer_min_kill_ts * 1000);
            data.execute_ts_string = execute_ts.getHours() + ':' + execute_ts.getMinutes() + ' on ' + execute_ts.getMonth() + '/' + execute_ts.getDate();
            data.min_kill_ts_string = min_kill_ts.getHours() + ':' + min_kill_ts.getMinutes() + ' on ' + min_kill_ts.getMonth() + '/' + min_kill_ts.getDate();

            this.$el.html(this.template(data));
            return this;
        }
    });
})();
