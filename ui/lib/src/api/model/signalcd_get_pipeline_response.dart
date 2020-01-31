part of swagger.api;

class SignalcdGetPipelineResponse {
  
  SignalcdPipeline pipeline = null;
  
  SignalcdGetPipelineResponse();

  @override
  String toString() {
    return 'SignalcdGetPipelineResponse[pipeline=$pipeline, ]';
  }

  SignalcdGetPipelineResponse.fromJson(Map<String, dynamic> json) {
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

  static List<SignalcdGetPipelineResponse> listFromJson(List<dynamic> json) {
    return json == null ? new List<SignalcdGetPipelineResponse>() : json.map((value) => new SignalcdGetPipelineResponse.fromJson(value)).toList();
  }

  static Map<String, SignalcdGetPipelineResponse> mapFromJson(Map<String, Map<String, dynamic>> json) {
    var map = new Map<String, SignalcdGetPipelineResponse>();
    if (json != null && json.length > 0) {
      json.forEach((String key, Map<String, dynamic> value) => map[key] = new SignalcdGetPipelineResponse.fromJson(value));
    }
    return map;
  }
}

