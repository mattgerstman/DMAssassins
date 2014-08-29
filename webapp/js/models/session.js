var app = app || {Models:{}, Views:{}, Routers:{}, Running:{}, Session:{}};

(function(){

	app.Models.Session = Backbone.Model.extend({
		
		url : config.WEB_ROOT+'session/',

		initialize : function(){
			//Ajax Request Configuration
			//To Set The CSRF Token To Request Header
/*
			$.ajaxSetup({
				headers : {
					'X-CSRF-Token' : csrf
				}
			});
*/

			//Check for sessionStorage support
			if(Storage && sessionStorage){
				this.supportStorage = true;
			}
		},

		get : function(key){
			if(this.supportStorage){
				var data = sessionStorage.getItem(key);
				if(data && data[0] === '{'){
					return JSON.parse(data);
				}else{
					return data;
				}
			}else{
				return Backbone.Model.prototype.get.call(this, key);
			}
		},


		set : function(key, value){
			if(this.supportStorage){
				sessionStorage.setItem(key, value);
			}else{
				Backbone.Model.prototype.set.call(this, key, value);
			}
			return this;
		},

		unset : function(key){
			if(this.supportStorage){
				sessionStorage.removeItem(key);
			}else{
				Backbone.Model.prototype.unset.call(this, key);
			}
			return this;	
		},

		clear : function(){
			if(this.supportStorage){
				sessionStorage.clear();  
			}else{
				Backbone.Model.prototype.clear(this);
			}
		},
		login: function(){

			var parent = this;

			FB.getLoginStatus(function(response){
				  if (response.status === 'connected') {
				    // Logged into your app and Facebook.

					console.log(response);
					parent.createSession(response);
					
//					var authKey = Base64.encode(userID+':'+token);


				  } else if (response.status === 'not_authorized') {
				    // The person is logged into Facebook, but not your app.
				    FB.login(function(response){
    					parent.createSession(response);
				    }, {scope:'public_profile,email,user_friends,user_photos'})
				    
				  } else {

					    FB.login(function(response){
	    					parent.createSession(response);
					    }, {scope:'public_profile,email,user_friends,user_photos'})

				    // The person is not logged into Facebook, so we're not sure if
				    // they are logged into this app or not.
				  }
			})

		},
		createSession: function(response, callback) {
			var data =  {
				'facebook_id': response.authResponse.userID,
				'facebook_token' : response.authResponse.accessToken
			}
			var that = this;
			var login = $.ajax({
				url : this.url,
				data : data,
				type : 'POST'
			});
			login.done(function(response){
				that.set('authenticated', true);

				that.set('token', JSON.stringify(response.response.token));
				that.set('user', JSON.stringify(response.response.user));
				that.setGame(response.response.game);
				that.storeBasicAuth(response.response)
				
				app.Running.ProfileModel = new app.Models.Profile(that.get('user'))
				
				if(that.get('redirectFrom')){
					var path = that.get('redirectFrom');
					that.unset('redirectFrom');
					Backbone.history.navigate(path, { trigger : true });
				}else{
					Backbone.history.navigate('', { trigger : true });
				}
			});
			login.fail(function(){
				Backbone.history.navigate('login', { trigger : true });
			});
			

		},
		setGame: function(game) {
			if (!game)
			{
				return;
			}
			this.set('game_id', game.game_id)
			this.set('game', JSON.stringify(game));
		},
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


		getAuth : function(callback){

			this.login()
/*
			var that = this;
			var Session = this.fetch();
			
			

			Session.done(function(response){
				that.set('authenticated', true);
				that.set('user', JSON.stringify(response.user));
			});

			Session.fail(function(response){
				response = JSON.parse(response.responseText);
				that.clear();
				csrf = response.csrf !== csrf ? response.csrf : csrf;
				that.initialize();
			});
*/

			Session.always(callback);
		},
		storeBasicAuth : function(data) {
			
			var user_id = data.user.user_id;
			this.set('user_id', user_id)

			var token = data.token
			var plainKey = user_id + ":" + token
			var base64Key = Base64.encode(plainKey);
			this.set('authKey', base64Key);
			this.setAuthHeader();	
		},
		setAuthHeader: function(){
			var base64Key = this.get('authKey');
			$.ajaxSetup({
				headers: { 'Authorization': "Basic " + base64Key }
			});

		}
	});
})()