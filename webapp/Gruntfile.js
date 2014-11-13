/*global module:false*/
module.exports = function(grunt) {

  // Project configuration.
  grunt.initConfig({
    // Metadata.
    pkg: grunt.file.readJSON('package.json'),
    banner: '/*!\n * <%= pkg.title || pkg.name %> - v<%= pkg.version %> - <%= pkg.code_name %> - ' +
      '<%= grunt.template.today("yyyy-mm-dd") %>\n' +
      '<%= pkg.homepage ? " * " + pkg.homepage + "\\n" : "" %>' +
      ' * Copyright (c) <%= grunt.template.today("yyyy") %> <%= pkg.author.name %>;\n' +
      ' */\n',
    // Task configuration.
    concat: {
      options: {
        banner: '<%= banner %>',
        stripBanners: true
      },
      dist: {
        src: [
          'js/config.js',
          'js/lib/*.js',      
          'js/models/*.js',
          'js/collections/*.js',
          'js/views/*.js',
          'js/routers/*.js',
          'js/*.js'
          ],
        dest: 'dist/<%= pkg.version %>/<%= pkg.name %>.js'
      }
    },
    uglify: {
      options: {
        banner: '<%= banner %>',
        sourceMap: true
      },
      dist: {
        src: '<%= concat.dist.dest %>',
        dest: 'dist/<%= pkg.version %>/<%= pkg.name %>.min.js'
      }
    },
    jshint: {
      files: [
          'js/config.js',
          'js/lib/*.js',
          'js/models/*.js',
          'js/collections/*.js',
          'js/views/*.js',
          'js/routers/*.js',
          'js/*.js'
        ],
      gruntfile: {
        src: 'Gruntfile.js'
      },      
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
          src : 'index.html.template',
          dest : 'index.html'
      },
      prod : {
          src : 'index.html.template',
          dest : 'index.html'
      }
    },
    less: {
      dev: {
        options: {
          banner: '<%= banner %>',          
        },
        files: {
          'dist/<%= pkg.version %>/DMAssassins.css': 'assets/styles/*.less'          
        }
      },
      prod: {
        options: {
          banner: '<%= banner %>',
          cleancss: true          
        },
        files: {          
          'dist/<%= pkg.version %>/DMAssassins.min.css': 'assets/styles/*.less'          
        }
      }
    },
    watch : {
      css: {
        files: 'assets/styles/*.less',
        tasks: ['less:dev']
      },
      js: {
        files: 'js/*/*.js',
        tasks: ['jshint']
      },
      index: {
        files: 'index.html.template',
        tasks: ['env:dev', 'preprocess:dev']
      }
      
    }
  });

  // These plugins provide necessary tasks.
  
  grunt.loadNpmTasks('grunt-contrib-concat');  
  grunt.loadNpmTasks('grunt-contrib-less');
  grunt.loadNpmTasks('grunt-contrib-uglify');
  grunt.loadNpmTasks('grunt-contrib-watch');
  grunt.loadNpmTasks('grunt-env');
  grunt.loadNpmTasks('grunt-preprocess');
  grunt.loadNpmTasks('grunt-contrib-jshint');

  // Default task.
  grunt.registerTask('dev', ['jshint', 'less:dev', 'env:dev', 'preprocess:dev']);
  grunt.registerTask('prod', ['concat', 'uglify', 'less:prod', 'env:prod', 'preprocess:prod']);

};
