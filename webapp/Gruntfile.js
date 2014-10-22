/*global module:false*/
module.exports = function(grunt) {

  // Project configuration.
  grunt.initConfig({
    // Metadata.
    pkg: grunt.file.readJSON('package.json'),
    banner: '/*! <%= pkg.title || pkg.name %> - v<%= pkg.version %> - ' +
      '<%= grunt.template.today("yyyy-mm-dd") %>\n' +
      '<%= pkg.homepage ? "* " + pkg.homepage + "\\n" : "" %>' +
      '* Copyright (c) <%= grunt.template.today("yyyy") %> <%= pkg.author.name %>;' +
      ' Licensed <%= _.pluck(pkg.licenses, "type").join(", ") %> */\n',
    // Task configuration.
    concat: {
      options: {
        banner: '<%= banner %>',
        stripBanners: true
      },
      dist: {
        src: [
          'js/lib/*.js',
          'js/models/*.js',
          'js/collections/*.js',
          'js/views/*.js',
          'js/routers/*.js',
          'js/*.js'
          ],
        dest: 'dist/<%= pkg.name %>.js'
      }
    },
    uglify: {
      options: {
        banner: '<%= banner %>'
      },
      dist: {
        src: '<%= concat.dist.dest %>',
        dest: 'dist/<%= pkg.version %>/<%= pkg.name %>.min.js'
      }
    },
    jshint: {
      files: '<%= concat.dist.src %>',
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
          BETA: '<%= pkg.beta %>'
      },
      prod : {
          NODE_ENV: 'PRODUCTION',
          BETA: '<%= pkg.beta %>',
      }
    },
    preprocess: {
      dev : {
          src : 'index.html.template',
          dest : 'index.html'
      },
      prod : {
          src : 'index.html.template',
          dest : 'index.html',
          options : {
              context : {
                  name : '<%= pkg.name %>',
                  version : '<%= pkg.version %>',                  
                  now : '<%= now %>',
                  ver : '<%= ver %>'
              }
          }
      }
    },
    cssmin: {
      add_banner: {
        options: {
          banner: '<%= banner %>'
        },
        files: {
            'dist/<%= pkg.version %>/DMAssassins.min.css': ['assets/styles/*.css']
        }
      }
    }
  });

  // These plugins provide necessary tasks.
  grunt.loadNpmTasks('grunt-contrib-concat');
  grunt.loadNpmTasks('grunt-contrib-uglify');
  grunt.loadNpmTasks('grunt-contrib-jshint');
  grunt.loadNpmTasks('grunt-preprocess');
  grunt.loadNpmTasks('grunt-env');
  grunt.loadNpmTasks('grunt-contrib-cssmin');

  // Default task.
  grunt.registerTask('dev', ['jshint', 'env:dev', 'preprocess:dev']);
  grunt.registerTask('prod', ['concat', 'uglify', 'cssmin', 'env:prod', 'preprocess:prod']);

};
