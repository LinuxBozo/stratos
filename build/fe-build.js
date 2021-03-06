/**
 * Gulp build file for Angular 2 Front End UI Code
 */
/* eslint-disable angular/log,no-console,no-process-env,angular/json-functions,no-sync */
(function () {
  'use strict';

  var gulp = require('gulp');
  // var _ = require('lodash');
  var del = require('delete');
  var spawn = require('child_process').spawn;
  var path = require('path');
  var os = require('os');

  var config = require('./gulp.config');
  var paths = config.paths;

  // Clean dist dir
  gulp.task('clean', function (next) {
    del(paths.dist + '**/*', {
      force: true
    }, next);
  });

  // Legacy task name
  gulp.task('clean:dist', gulp.series('clean'));

  gulp.task('ng-build', function (cb) {
    var rootFolder = path.resolve(__dirname, '..');
    var cmd = 'npm';
    var windowsEnvironment = os.platform().startsWith('win');
    if (windowsEnvironment) {
      cmd = 'npm.cmd';
    }
    var child = spawn(cmd, ['run', 'build'], {
      cwd: rootFolder
    });
    child.stdout.on('data', function (data) {
      console.log(data.toString());
    });
    child.stderr.on('data', function (data) {
      console.log(data.toString());
    });
    child.on('error', function (err) {
      console.log(err);
      cb(err);
    });
    child.on('close', function (code) {
      var err = code === 0 ? undefined : 'Build exited with code: ' + code;
      cb(err);
    });
  });

  // Production build
  gulp.task('build', gulp.series(
    'clean',
    'ng-build'
  ));

  // Default task is to build for production
  gulp.task('default', gulp.series('build'));

})();
