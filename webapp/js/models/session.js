// Manages all local storage information and helps keep various models in sync
// js/models/session

var app = app || {Models:{}, Views:{}, Routers:{}, Running:{}, Session:{}};

(function(){

	app.Models.Session = Backbone.Model.extend({
		
		url : config.WEB_ROOT+'session/',

		initialize : function(){
		
			// Check for localstorage support
			if(Storage && sessionStorage){
				this.supportStorage = true;
			}
		},

		// returns data stored in the session
		get : function(key) {
			if (this.supportStorage) {
				var data = sessionStorage.getItem(key);
				if(data && data[0] === '{'){
					return JSON.parse(data);
				} else {
					return data;
				}
			} else {
				return Backbone.Model.prototype.get.call(this, key);
			}
		},

		// sets a session variable
		set : function(key, value) {
			if(this.supportStorage) {
				sessionStorage.setItem(key, value);
			} else {
				Backbone.Model.prototype.set.call(this, key, value);
			}
			return this;
		},

		// unsets a session 
		unset : function(key){
			if(this.supportStorage){
				sessionStorage.removeItem(key);
			}else{
				Backbone.Model.prototype.unset.call(this, key);
			}
			return this;	
		},

		// clears all data from the session
		clear : function(){
			if(this.supportStorage) {
				sessionStorage.clear();  
			} else {
				Backbone.Model.prototype.clear(this);
			}
		},
		
		// calls the facebook login function and handles it appropriately
		// if they are logged into facebook and connected to the app a session is created automatically
		// otherwise a popup will appear and handle the session situation
		login: function(){

			var parent = this;

			FB.getLoginStatus(function(response){
				  if (response.status === 'connected') {
				    // Logged into your app and Facebook.
					//console.log(response);
					parent.createSession(response);


				  } else if (response.status === 'not_authorized') {

				    // The person is logged into Facebook, but not your app.
				    FB.login(function(response){
    					parent.createSession(response);
    					
    					// scope are the facebook permissions we're requesting 
				    }, {scope:'public_profile,email,user_friends,user_photos'})
				    
				  } else {

					    FB.login(function(response){
	    					parent.createSession(response);
	    					
	    					// scope are the facebook permissions we're requesting
					    }, {scope:'public_profile,email,user_friends,user_photos'})

				    // The person is not logged into Facebook, so we're not sure if
				    // they are logged into this app or not.
				  }
			})

		},
		
		// takes a facebook response and creates a session from it
		createSession: function(response) {

			var data =  {
				'facebook_id': response.authResponse.userID,
				'facebook_token' : response.authResponse.accessToken
			}
			
			var that = this;
			
			// performs the ajax request to the server to get session data
			var login = $.ajax({
				url : this.url,
				data : data,
				type : 'POST'
			});
			
			// after the ajax request run this function
			login.done(function(response){

				// store a goto boolean to determine if we're authenticated
				that.set('authenticated', true);

				// store the user in a session, this is game agnostic so it especially fits here
				that.set('user', JSON.stringify(response.response.user));
				
				// store the current game in the session
				that.set('game', JSON.stringify(response.response.game))
				
				// store the basic auth token in the session in case we need to reload it on app launch
				that.storeBasicAuth(response.response)
				
				// load the profile model for the user
				app.Running.ProfileModel = new app.Models.Profile(that.get('user'))
				app.Running.TargetModel = new app.Models.Target({assassin_id: that.get('user').user_id})
				if (app.Running.NavGameView !== undefined) {
					app.Running.NavGameView.render();
				}
				
				Backbone.history.navigate(Backbone.history.fragment, { trigger : true });
				
			});
			
			// if theres a login error direct them to the login screen
			login.fail(function(){
				Backbone.history.navigate('login', { trigger : true });
			});
			

		},
		// clear all the session data and post it to the server
		logout : function(callback){
			var that = this;
			$.ajax({
				url : this.url + '/logout',
				type : 'DELETE'
			}).done(function(response){
				//Clear all session data
				that.clear();
				//Set the new csrf token to csrf vaiable and
				//call initialize to update the $.ajaxSetup 
				// with new csrf
				csrf = response.csrf;
				that.initialize();
				callback();
			});
		},
		
		// stores all the basic auth variables in the session
		storeBasicAuth : function(data) {
			
			var user_id = data.user.user_id;
			this.set('user_id', user_id)

			var token = data.token
			var plainKey = user_id + ":" + token
			var base64Key = Base64.encode(plainKey);
			this.set('authKey', base64Key);
			this.setAuthHeader();	
		},
		
		// sets the Basic Auth header for all ajax requests
		setAuthHeader: function(){
			var base64Key = this.get('authKey');
			$.ajaxSetup({
				headers: { 'Authorization': "Basic " + base64Key }
			});

		}
	});
})()