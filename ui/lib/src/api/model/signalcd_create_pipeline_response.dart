part of swagger.api;

class SignalcdCreatePipelineResponse {
  
  SignalcdPipeline pipeline = null;
  
  SignalcdCreatePipelineResponse();

  @override
  String toString() {
    return 'SignalcdCreatePipelineResponse[pipeline=$pipeline, ]';
  }

  SignalcdCreatePipelineResponse.fromJson(Map<String, dynamic> json) {
    if (json == null) return;
    pipeline =
      
      
      new SignalcdPipeline.fromJson(json['pipeline'])
;
  }

  Map<String, dynamic> toJson() {
    return {
      'pipeline': pipeline
     };
  }

  static List<SignalcdCreatePipelineResponse> listFromJson(List<dynamic> json) {
    return json == null ? new List<SignalcdCreatePipelineResponse>() : json.map((value) => new SignalcdCreatePipelineResponse.fromJson(value)).toList();
  }

  static Map<String, SignalcdCreatePipelineResponse> mapFromJson(Map<String, Map<String, dynamic>> json) {
    var map = new Map<String, SignalcdCreatePipelineResponse>();
    if (json != null && json.length > 0) {
      json.forEach((String key, Map<String, dynamic> value) => map[key] = new SignalcdCreatePipelineResponse.fromJson(value));
    }
    return map;
  }
}

