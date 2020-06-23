/**
 * SignalCD
 * No description provided (generated by Openapi Generator https://github.com/openapitools/openapi-generator)
 *
 * The version of the OpenAPI document: 0.0.0
 * 
 *
 * NOTE: This class is auto generated by OpenAPI Generator (https://openapi-generator.tech).
 * https://openapi-generator.tech
 * Do not edit the class manually.
 *
 */

(function(root, factory) {
  if (typeof define === 'function' && define.amd) {
    // AMD.
    define(['expect.js', process.cwd()+'/src/index'], factory);
  } else if (typeof module === 'object' && module.exports) {
    // CommonJS-like environments that support module.exports, like Node.
    factory(require('expect.js'), require(process.cwd()+'/src/index'));
  } else {
    // Browser globals (root is window)
    factory(root.expect, root.SignalCd);
  }
}(this, function(expect, SignalCd) {
  'use strict';

  var instance;

  beforeEach(function() {
    instance = new SignalCd.PipelineSteps();
  });

  var getProperty = function(object, getter, property) {
    // Use getter method if present; otherwise, get the property directly.
    if (typeof object[getter] === 'function')
      return object[getter]();
    else
      return object[property];
  }

  var setProperty = function(object, setter, property, value) {
    // Use setter method if present; otherwise, set the property directly.
    if (typeof object[setter] === 'function')
      object[setter](value);
    else
      object[property] = value;
  }

  describe('PipelineSteps', function() {
    it('should create an instance of PipelineSteps', function() {
      // uncomment below and update the code to test PipelineSteps
      //var instane = new SignalCd.PipelineSteps();
      //expect(instance).to.be.a(SignalCd.PipelineSteps);
    });

    it('should have the property name (base name: "name")', function() {
      // uncomment below and update the code to test the property name
      //var instane = new SignalCd.PipelineSteps();
      //expect(instance).to.be();
    });

    it('should have the property image (base name: "image")', function() {
      // uncomment below and update the code to test the property image
      //var instane = new SignalCd.PipelineSteps();
      //expect(instance).to.be();
    });

    it('should have the property commands (base name: "commands")', function() {
      // uncomment below and update the code to test the property commands
      //var instane = new SignalCd.PipelineSteps();
      //expect(instance).to.be();
    });

  });

}));
