part of openapi.api;

class SignalcdGetPipelineResponse {
  
  SignalcdPipeline pipeline = null;
  SignalcdGetPipelineResponse();

  @override
  String toString() {
    return 'SignalcdGetPipelineResponse[pipeline=$pipeline, ]';
  }

  SignalcdGetPipelineResponse.fromJson(Map<String, dynamic> json) {
    if (json == null) return;
    pipeline = (json['pipeline'] == null) ?
      null :
      SignalcdPipeline.fromJson(json['pipeline']);
  }

  Map<String, dynamic> toJson() {
    Map <String, dynamic> json = {};
    if (pipeline != null)
      json['pipeline'] = pipeline;
    return json;
  }

  static List<SignalcdGetPipelineResponse> listFromJson(List<dynamic> json) {
    return json == null ? List<SignalcdGetPipelineResponse>() : json.map((value) => SignalcdGetPipelineResponse.fromJson(value)).toList();
  }

  static Map<String, SignalcdGetPipelineResponse> mapFromJson(Map<String, dynamic> json) {
    var map = Map<String, SignalcdGetPipelineResponse>();
    if (json != null && json.isNotEmpty) {
      json.forEach((String key, dynamic value) => map[key] = SignalcdGetPipelineResponse.fromJson(value));
    }
    return map;
  }

  // maps a json object with a list of SignalcdGetPipelineResponse-objects as value to a dart map
  static Map<String, List<SignalcdGetPipelineResponse>> mapListFromJson(Map<String, dynamic> json) {
    var map = Map<String, List<SignalcdGetPipelineResponse>>();
     if (json != null && json.isNotEmpty) {
       json.forEach((String key, dynamic value) {
         map[key] = SignalcdGetPipelineResponse.listFromJson(value);
       });
     }
     return map;
  }
}

