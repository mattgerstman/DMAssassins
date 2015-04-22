(function() {
  var template = Handlebars.template, templates = Handlebars.templates = Handlebars.templates || {};
templates['nav.template.html'] = template({"1":function(depth0,helpers,partials,data) {
    return "disabled\">";
},"3":function(depth0,helpers,partials,data) {
    var stack1;

  return "          <li class=\"js-dropdown-parent js-dropdown-parent-admin\">\n            <a data-toggle=\"dropdown\" class=\"dropdown-toggle\" href=\"#\">\n              <div class=\"hidden-sm\">\n                <span>\n                  Admin\n                </span>\n                <span class=\"caret\"></span>\n              </div>\n              <div class=\"visible-sm\">\n                <span>\n                  Admin\n                </span>\n                <span class=\"caret\"></span>\n              </div>\n            </a>\n            <ul class=\"dropdown-menu\" role=\"menu\" id=\"js-nav-admin-dropdown\">\n              <li dropdown=\"admin\" id=\"js-nav-users\"><a class=\"js-nav-link\" href=\"#users\">Manage Users</a></li>\n              <li dropdown=\"admin\" id=\"js-nav-edit-rules\"><a class=\"js-nav-link\" href=\"#edit_rules\">Edit Rules</a></li>\n              <li dropdown=\"admin\" id=\"js-nav-game-settings\"><a class=\"js-nav-link\" href=\"#game_settings\">Game Settings</a></li>\n              <!-- Super Admin Links -->"
    + ((stack1 = helpers['if'].call(depth0,(depth0 != null ? depth0.is_super_admin : depth0),{"name":"if","hash":{},"fn":this.program(4, data, 0),"inverse":this.noop,"data":data})) != null ? stack1 : "")
    + "</ul>\n          </li>\n          <!-- Captain Links -->";
},"4":function(depth0,helpers,partials,data) {
    return "                  <li dropdown=\"admin\" id=\"js-nav-targets\"><a class=\"js-nav-link\" href=\"#targets\">View Targets</a></li>";
},"6":function(depth0,helpers,partials,data) {
    return "            <li id=\"js-nav-users\"><a class=\"js-nav-link\" href=\"#users\">Manage Users</a></li>";
},"compiler":[6,">= 2.0.0-beta.1"],"main":function(depth0,helpers,partials,data) {
    var stack1;

  return "<nav role=\"navigation\" class=\"navbar navbar-default navbar-fixed-top\">\n  <div class=\"container\">\n    <!-- Brand and toggle get grouped for better mobile display -->\n    <div class=\"navbar-header\">\n      <button data-target=\"#main_nav\" data-toggle=\"collapse\" class=\"navbar-toggle\" type=\"button\">\n        <span class=\"sr-only\">Toggle navigation</span>\n        <span class=\"icon-bar\"></span>\n        <span class=\"icon-bar\"></span>\n        <span class=\"icon-bar\"></span>\n      </button>\n      <!-- @if NODE_ENV='PRODUCTION' -->\n      <a href=\"#\" class=\"navbar-brand\">DMAssassins</a>\n      <!-- @endif -->\n      <!-- @if NODE_ENV='DEVELOPMENT' -->\n      <a href=\"#\" class=\"navbar-brand\">DevAssassins</a>\n      <!-- @endif -->\n    </div>\n\n    <!-- Collect the nav links, forms, and other content for toggling -->\n    <div class=\"collapse navbar-collapse\" id=\"main_nav\">\n      <!-- Left Navbar -->\n      <ul class=\"nav navbar-nav\">\n        <li class=\"js-nav-target\n"
    + ((stack1 = helpers['if'].call(depth0,(depth0 != null ? depth0.target_disabled : depth0),{"name":"if","hash":{},"fn":this.program(1, data, 0),"inverse":this.noop,"data":data})) != null ? stack1 : "")
    + "<a class=\"js-nav-link\" href=\"#target\">Target</a>\n        </li>\n        <li class=\"js-nav-my-profile\"><a class=\"js-nav-link\" href=\"#my_profile\">My Profile</a></li>\n        <li class=\"js-nav-leaderboard\"><a class=\"js-nav-link\" href=\"#leaderboard\">Leaderboard</a></li>\n        <li class=\"js-nav-rules\"><a class=\"js-nav-link\" href=\"#rules\">Rules</a></li>\n      </ul>\n\n      <!-- Right Navbar -->\n      <ul class=\"nav navbar-nav navbar-right\">\n        <li class=\"js-dropdown-parent js-dropdown-parent-games\">\n\n        </li>\n        <!-- Admin Links -->\n"
    + ((stack1 = helpers['if'].call(depth0,(depth0 != null ? depth0.is_admin : depth0),{"name":"if","hash":{},"fn":this.program(3, data, 0),"inverse":this.noop,"data":data})) != null ? stack1 : "")
    + ((stack1 = helpers['if'].call(depth0,(depth0 != null ? depth0.is_captain : depth0),{"name":"if","hash":{},"fn":this.program(6, data, 0),"inverse":this.noop,"data":data})) != null ? stack1 : "")
    + "<li id=\"js-nav-logout\"><a class=\"js-nav-link\" href=\"#logout\">Logout</a></li>\n        </ul>\n    </div> <!-- /.navbar-collapse -->\n  </div> <!-- /.container -->\n</nav> <!-- /.navbar -->\n";
},"useData":true});
})();