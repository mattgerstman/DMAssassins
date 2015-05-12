// Used to handle user permissions
// All APIs should be equipped with the right set of permissions so this simply prevents users from stumbling onto a page they can't use

var AuthUtils = {
    // Returns a map of role data
    getRolesMap: function() {
        var rolesMap = {
            dm_user:        {value: 0, pretty_name: "User"},
            dm_captain:     {value: 1, pretty_name: "Captain"},
            dm_admin:       {value: 2, pretty_name: "Admin"},
            dm_super_admin: {value: 3, pretty_name: "Super Admin"},
        };
        return rolesMap;
    },
    // Returns the role map a given role has access to
    getRolesMapFor: function(role, teams_enabled) {
        if (!role) {
            return {};
        }
        var filteredMap = {};
        var rolesMap = this.getRolesMap();
        for (var key in rolesMap)  {
            if (teams_enabled === false) {
                // Strip out captains if teams arent enabled
                if (rolesMap[key].value === 1)
                    continue;
            }
            // Compare the given role and the role to return
            if (rolesMap[key].value <= rolesMap[role].value)
                filteredMap[key] = rolesMap[key];
        }
        return filteredMap;
    },
    // Gets the pretty name for a role
    getRolePrettyName: function(role) {
        var rolesMap = this.getRolesMap();
        return rolesMap[role].pretty_name;
    },
    // Compare two roles
    isRoleAllowed: function(userRole, minRoleForAccess) {
        var rolesMap = this.getRolesMap();
        if ((userRole === undefined) || (rolesMap[userRole] === undefined))
            return undefined;

        return rolesMap[userRole].value >= rolesMap[minRoleForAccess].value;
    },
    // Determine if the given role is a user or above
    requiresUser: function(userRole) {
        return this.isRoleAllowed(userRole, 'dm_user');
    },
    // Determine if the given role is a captain or above
    requiresCaptain: function(userRole) {
        return this.isRoleAllowed(userRole, 'dm_captain');
    },
        // Determine if the given role is an admin or above
    requiresAdmin: function(userRole) {
        return this.isRoleAllowed(userRole, 'dm_admin');
    },
    // Determine if the given role is a super admin
    requiresSuperAdmin: function(userRole) {
        return this.isRoleAllowed(userRole, 'dm_super_admin');
    }
};
