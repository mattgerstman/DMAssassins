/*global module:false*/
module.exports = function(grunt) {

  // MODIFIED: add require for connect-modewrite
  var modRewrite = require('connect-modrewrite');
  var url = require("url");
  var fs = require("fs");
  var path = require("path");

  // Project configuration.
  grunt.initConfig({
    // Metadata.
    pkg: grunt.file.readJSON('package.json'),
    banner: '/*!\n * <%= pkg.title || pkg.name %> - v<%= pkg.version %> - <%= pkg.code_name %> - ' +
      '<%= grunt.template.today("yyyy-mm-dd") %>\n' +
      '<%= pkg.homepage ? " * " + pkg.homepage + "\\n" : "" %>' +
      ' * Copyright (c) <%= grunt.template.today("yyyy") %> <%= pkg.author.name %>;\n' +
      ' */',
    // Task configuration.
    dependencies: {
      js: [
        'js/config.js',
        'js/lib/*.js',
        'js/models/*.js',
        'js/collections/*.js',
        'dist/<%= pkg.version %>/templates/user.js',
        'js/views/user/*.js',
        'js/routers/*.js',
        'js/*.js'
      ],
      captain: [
        'dist/captain/<%= pkg.version %>/templates/*.js',
        'js/views/captain/*.js',
      ],
      admin: [
        'dist/admin/<%= pkg.version %>/templates/*.js',
        'js/views/admin/*.js',
      ],
      superadmin: [
        'dist/superadmin/<%= pkg.version %>/templates/*.js',
        'js/views/superadmin/*.js',
      ]
    },
    jst: {
      site: {
        options: {
          prettify: true,
          namespace: 'app.Templates',
          processName: function(path) {
            return path.substring(path.lastIndexOf('/')+1, path.lastIndexOf('.'));
          }
        },
        files: {
          "dist/<%= pkg.version %>/templates/user.js": ["templates/user/*.html"],
          "dist/captain/<%= pkg.version %>/templates/captain.js": ["templates/captain/*.html"],
          "dist/admin/<%= pkg.version %>/templates/admin.js": ["templates/admin/*.html"],
          "dist/superadmin/<%= pkg.version %>/templates/superadmin.js": ["templates/superadmin/*.html"]
        }
      },
      plotTwists: {
        options: {
          prettify: true,
          namespace: 'app.Templates.PlotTwist',
          processName: function(path) {
            return path.substring(path.lastIndexOf('/')+1, path.lastIndexOf('.'));
          }
        },
        files: {
          "dist/admin/<%= pkg.version %>/templates/plot-twists.js": ["templates/plot-twists/*.html"],
        }
      }
    },
    uglify: {
      options: {
        banner: '<%= banner %>',
        sourceMap: true
      },
      dist: {
        src: [
          '<%= dependencies.js %>'
        ],
        dest: 'dist/<%= pkg.version %>/<%= pkg.name %>.min.js'
      },
      captain: {
        src: [
          '<%= dependencies.captain %>'
        ],
        dest: 'dist/captain/<%= pkg.version %>/<%= pkg.name %>-captain.min.js'
      },
      admin: {
        src: [
          '<%= dependencies.admin %>'
        ],
        dest: 'dist/admin/<%= pkg.version %>/<%= pkg.name %>-admin.min.js'
      },
      superadmin: {
        src: [
          '<%= dependencies.superadmin %>'
        ],
        dest: 'dist/superadmin/<%= pkg.version %>/<%= pkg.name %>-superadmin.min.js'
      }
    },
    jshint: {
      files: [
        '<%= dependencies.js %>',
        '!dist/<%= pkg.version %>/templates/*.js'
      ],
      gruntfile: {
        src: 'Gruntfile.js'
      },
    },
    lintspaces: {
	    javascript: {
        src: [
          '<%= dependencies.js %>',
          '!dist/<%= pkg.version %>/templates/*.js'
        ],
        options: {
          newline: true,
          newlineMaximum: 2,
          trailingspaces: true,
          indentation: 'spaces',
          spaces: 4,
          ignores: ['js-comments']
        }
	    },
	    grunt: {
  	    src: [
    	    'Gruntfile.js'
  	    ],
  	    options: {
          newline: true,
          newlineMaximum: 2,
          trailingspaces: true,
          indentation: 'spaces',
          spaces: 2,
          ignores: ['js-comments']
  	    }
	    }
	  },
    env : {
      options : {
        VERSION: '<%= pkg.version %>',
      },
      dev: {
        NODE_ENV: 'DEVELOPMENT',
        BETA: '<%= pkg.beta %>',
        BANNER: '<%= banner %>',
      },
      prod : {
        NODE_ENV: 'PRODUCTION',
        BETA: '<%= pkg.beta %>',
        BANNER: '<%= banner %>',
      }
    },
    preprocess: {
      dev : {
        files: {
          'index.html' : 'index.html.template'
        }
      },
      prod : {
        files: {
          'index.html' : 'index.html.template'
        }
      }
    },
    injector: {
      options: {

      },
      dev: {
        files: {
          'index.html' : [
            '<%= dependencies.js %>',
            '<%= dependencies.captain %>',
            '<%= dependencies.admin %>',
            '<%= dependencies.superadmin %>',
            '<%= less.dev.dest %>'
          ]
        }
      },
      prod: {
        files: {
          'index.html' : [
            '<%= uglify.dist.dest %>',
            '<%= less.prod.dest %>'
          ]
        }
      }
    },
    less: {
      dev: {
        options: {
          banner: '<%= banner %>',
          strictMath: true,
          sourceMap: true,
          outputSourceFiles: true,
          sourceMapURL: 'DMAssassins.css.map',
          sourceMapFilename: 'dist/<%= pkg.version %>/DMAssassins.css.map'
        },
        src: 'assets/styles/DMAssassins.less',
        dest: 'dist/<%= pkg.version %>/DMAssassins.css'
      },
      prod: {
        options: {
          banner: '<%= banner %>',
          strictMath: true,
          sourceMap: true,
          outputSourceFiles: true,
          sourceMapURL: 'DMAssassins.css.map',
          sourceMapFilename: 'dist/<%= pkg.version %>/DMAssassins.css.map',
          cleancss: true
        },
        src: 'assets/styles/DMAssassins.less',
        dest: 'dist/<%= pkg.version %>/DMAssassins.min.css'

      }
    },
    watch : {
      css: {
        files: 'assets/styles/*.less',
        tasks: ['less:dev']
      },
      js: {
        files: 'js/*/*/*.js',
        tasks: ['lintspaces', 'jshint']
      },
      index: {
        files: 'index.html.template',
        tasks: ['env:dev', 'preprocess:dev', 'injector:dev']
      },
      templates: {
        files: 'templates/*/*',
        tasks: ['jst']
      }
    },
    browserSync: {
      bsFiles: {
        src : [
          'js/*.js',
          'js/*/*.js',
          'js/*/*/*.js',
          'js/*/*/*/*.js',
          'index.html',
          'dist/*.css'
        ]
      },
      options: {
        port: 8888,
        open: 'ui',
        ui: {
          port: 9999
        },
        server: {
          baseDir: "./",
          target:"http://assassins.com",
          middleware: function(req, res, next) {
            var fileName = url.parse(req.url);
            fileName = fileName.href.split(fileName.search).join("");
            var fileExists = fs.existsSync(__dirname + fileName);
            if (!fileExists && fileName.indexOf("browser-sync-client") < 0) {
                req.url = "/index.html";
            }
            return next();
          }
        },
      }
    }
  });

  // These plugins provide necessary tasks.
  grunt.loadNpmTasks('grunt-browser-sync');
  grunt.loadNpmTasks('grunt-contrib-jshint');
  grunt.loadNpmTasks('grunt-contrib-jst');
  grunt.loadNpmTasks('grunt-contrib-less');
  grunt.loadNpmTasks('grunt-injector');
  grunt.loadNpmTasks('grunt-lintspaces');
  grunt.loadNpmTasks('grunt-contrib-uglify');
  grunt.loadNpmTasks('grunt-contrib-watch');
  grunt.loadNpmTasks('grunt-env');
  grunt.loadNpmTasks('grunt-preprocess');

  // Default task.
  grunt.registerTask('dev', ["jst", 'lintspaces', 'jshint', 'less:dev', 'env:dev', 'preprocess:dev', 'injector:dev']);
  grunt.registerTask('prod', ["jst", 'uglify', 'less:prod', 'env:prod', 'preprocess:prod', 'injector:prod']);
  grunt.registerTask('server', 'browserSync');
  grunt.registerTask('default', ['dev']);
};
