part of swagger.api;

class SignalcdListPipelinesResponse {
  
  List<SignalcdPipeline> pipelines = [];
  
  SignalcdListPipelinesResponse();

  @override
  String toString() {
    return 'SignalcdListPipelinesResponse[pipelines=$pipelines, ]';
  }

  SignalcdListPipelinesResponse.fromJson(Map<String, dynamic> json) {
    if (json == null) return;
    pipelines =
      SignalcdPipeline.listFromJson(json['pipelines'])
;
  }

  Map<String, dynamic> toJson() {
    return {
      'pipelines': pipelines
     };
  }

  static List<SignalcdListPipelinesResponse> listFromJson(List<dynamic> json) {
    return json == null ? new List<SignalcdListPipelinesResponse>() : json.map((value) => new SignalcdListPipelinesResponse.fromJson(value)).toList();
  }

  static Map<String, SignalcdListPipelinesResponse> mapFromJson(Map<String, Map<String, dynamic>> json) {
    var map = new Map<String, SignalcdListPipelinesResponse>();
    if (json != null && json.length > 0) {
      json.forEach((String key, Map<String, dynamic> value) => map[key] = new SignalcdListPipelinesResponse.fromJson(value));
    }
    return map;
  }
}

