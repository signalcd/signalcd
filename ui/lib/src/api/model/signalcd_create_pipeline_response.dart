part of openapi.api;

class SignalcdCreatePipelineResponse {
  
  SignalcdPipeline pipeline = null;
  SignalcdCreatePipelineResponse();

  @override
  String toString() {
    return 'SignalcdCreatePipelineResponse[pipeline=$pipeline, ]';
  }

  SignalcdCreatePipelineResponse.fromJson(Map<String, dynamic> json) {
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

  static List<SignalcdCreatePipelineResponse> listFromJson(List<dynamic> json) {
    return json == null ? List<SignalcdCreatePipelineResponse>() : json.map((value) => SignalcdCreatePipelineResponse.fromJson(value)).toList();
  }

  static Map<String, SignalcdCreatePipelineResponse> mapFromJson(Map<String, dynamic> json) {
    var map = Map<String, SignalcdCreatePipelineResponse>();
    if (json != null && json.isNotEmpty) {
      json.forEach((String key, dynamic value) => map[key] = SignalcdCreatePipelineResponse.fromJson(value));
    }
    return map;
  }

  // maps a json object with a list of SignalcdCreatePipelineResponse-objects as value to a dart map
  static Map<String, List<SignalcdCreatePipelineResponse>> mapListFromJson(Map<String, dynamic> json) {
    var map = Map<String, List<SignalcdCreatePipelineResponse>>();
     if (json != null && json.isNotEmpty) {
       json.forEach((String key, dynamic value) {
         map[key] = SignalcdCreatePipelineResponse.listFromJson(value);
       });
     }
     return map;
  }
}

