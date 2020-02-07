part of openapi.api;

class SignalcdListPipelinesResponse {
  
  List<SignalcdPipeline> pipelines = [];
  SignalcdListPipelinesResponse();

  @override
  String toString() {
    return 'SignalcdListPipelinesResponse[pipelines=$pipelines, ]';
  }

  SignalcdListPipelinesResponse.fromJson(Map<String, dynamic> json) {
    if (json == null) return;
    pipelines = (json['pipelines'] == null) ?
      null :
      SignalcdPipeline.listFromJson(json['pipelines']);
  }

  Map<String, dynamic> toJson() {
    Map <String, dynamic> json = {};
    if (pipelines != null)
      json['pipelines'] = pipelines;
    return json;
  }

  static List<SignalcdListPipelinesResponse> listFromJson(List<dynamic> json) {
    return json == null ? List<SignalcdListPipelinesResponse>() : json.map((value) => SignalcdListPipelinesResponse.fromJson(value)).toList();
  }

  static Map<String, SignalcdListPipelinesResponse> mapFromJson(Map<String, dynamic> json) {
    var map = Map<String, SignalcdListPipelinesResponse>();
    if (json != null && json.isNotEmpty) {
      json.forEach((String key, dynamic value) => map[key] = SignalcdListPipelinesResponse.fromJson(value));
    }
    return map;
  }

  // maps a json object with a list of SignalcdListPipelinesResponse-objects as value to a dart map
  static Map<String, List<SignalcdListPipelinesResponse>> mapListFromJson(Map<String, dynamic> json) {
    var map = Map<String, List<SignalcdListPipelinesResponse>>();
     if (json != null && json.isNotEmpty) {
       json.forEach((String key, dynamic value) {
         map[key] = SignalcdListPipelinesResponse.listFromJson(value);
       });
     }
     return map;
  }
}

