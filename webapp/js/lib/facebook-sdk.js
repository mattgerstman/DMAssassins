$(function() {

    // This is called with the results from from FB.getLoginStatus().
    function statusChangeCallback(response) {
        // The response object is returned with a status field that lets the
        // app know the current login status of the person.
        // Full docs on the response object can be found in the documentation
        // for FB.getLoginStatus().

        if (!app)
            return;

        if (!app.Session)
            return;

        if (app.Session.get('authenticated') === true) {
            // If we're already authenticated recover the previous session
            app.Session.recoverSession(response);
            return;
        }
        return;
    }
    FB.init({
        appId      : config.APP_ID,
        cookie     : true,  // enable cookies to allow the server to access
                            // the session
        xfbml      : true,  // parse social plugins on this page
        version    : 'v2.2' // use version 2.2
    });

    app.Running.FB = FB;

    // Now that we've initialized the JavaScript SDK, we call
    // FB.getLoginStatus().  This function gets the state of the
    // person visiting this page and can return one of three states to
    // the callback you provide.  They can be:
    //
    // 1. Logged into your app ('connected')
    // 2. Logged into Facebook, but not your app ('not_authorized')
    // 3. Not logged into Facebook and can't tell if they are logged into
    //    your app or not.
    //
    // These three cases are handled in the callback function.


    FB.getLoginStatus(function(response) {
        statusChangeCallback(response);
    });
});
