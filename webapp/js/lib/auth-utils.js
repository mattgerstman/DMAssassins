// Used to handle user permissions
// All APIs should be equipped with the right set of permissions so this simply prevents users from stumbling onto a page they can't use

var AuthUtils = {
    getRolesMap: function() {
        var rolesMap = {
            dm_user:        {value: 0, pretty_name: "User"},
            dm_captain:     {value: 1, pretty_name: "Captain"},
            dm_admin:       {value: 2, pretty_name: "Admin"},
            dm_super_admin: {value: 3, pretty_name: "Super Admin"},
        }
        return rolesMap;    
    },
    getRolePrettyName: function(role) {
        var rolesMap = this.getRolesMap();
        return rolesMap[role]['pretty_name'];
    },
    isRoleAllowed: function(userRole, minRoleForAccess) {
        var rolesMap = this.getRolesMap();
        if ((userRole === undefined) || (rolesMap[userRole] === undefined))
            return undefined
        
        return rolesMap[userRole]['value'] >= rolesMap[minRoleForAccess]['value'];        
    },
    requiresUser: function(userRole) {
        return this.isRoleAllowed(userRole, 'dm_user');
    },
    requiresCaptain: function(userRole) {
        return this.isRoleAllowed(userRole, 'dm_captain');
    },
    requiresAdmin: function(userRole) {
        return this.isRoleAllowed(userRole, 'dm_admin');
    },
    requiresSuperAdmin: function(userRole) {
        return this.isRoleAllowed(userRole, 'dm_super_admin');
    }
}

