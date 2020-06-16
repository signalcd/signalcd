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

import ApiClient from '../ApiClient';

/**
 * The Pipeline model module.
 * @module model/Pipeline
 * @version 0.0.0
 */
class Pipeline {
    /**
     * Constructs a new <code>Pipeline</code>.
     * @alias module:model/Pipeline
     * @param id {String} 
     */
    constructor(id) { 
        
        Pipeline.initialize(this, id);
    }

    /**
     * Initializes the fields of this object.
     * This method is used by the constructors of any subclasses, in order to implement multiple inheritance (mix-ins).
     * Only for internal use.
     */
    static initialize(obj, id) { 
        obj['id'] = id;
    }

    /**
     * Constructs a <code>Pipeline</code> from a plain JavaScript object, optionally creating a new instance.
     * Copies all relevant properties from <code>data</code> to <code>obj</code> if supplied or a new instance if not.
     * @param {Object} data The plain JavaScript object bearing properties of interest.
     * @param {module:model/Pipeline} obj Optional instance to populate.
     * @return {module:model/Pipeline} The populated <code>Pipeline</code> instance.
     */
    static constructFromObject(data, obj) {
        if (data) {
            obj = obj || new Pipeline();

            if (data.hasOwnProperty('id')) {
                obj['id'] = ApiClient.convertToType(data['id'], 'String');
            }
            if (data.hasOwnProperty('name')) {
                obj['name'] = ApiClient.convertToType(data['name'], 'String');
            }
            if (data.hasOwnProperty('created')) {
                obj['created'] = ApiClient.convertToType(data['created'], 'Date');
            }
        }
        return obj;
    }


}

/**
 * @member {String} id
 */
Pipeline.prototype['id'] = undefined;

/**
 * @member {String} name
 */
Pipeline.prototype['name'] = undefined;

/**
 * @member {Date} created
 */
Pipeline.prototype['created'] = undefined;






export default Pipeline;
