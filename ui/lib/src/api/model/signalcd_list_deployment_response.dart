part of openapi.api;

class SignalcdListDeploymentResponse {
  
  List<SignalcdDeployment> deployments = [];
  SignalcdListDeploymentResponse();

  @override
  String toString() {
    return 'SignalcdListDeploymentResponse[deployments=$deployments, ]';
  }

  SignalcdListDeploymentResponse.fromJson(Map<String, dynamic> json) {
    if (json == null) return;
    deployments = (json['deployments'] == null) ?
      null :
      SignalcdDeployment.listFromJson(json['deployments']);
  }

  Map<String, dynamic> toJson() {
    Map <String, dynamic> json = {};
    if (deployments != null)
      json['deployments'] = deployments;
    return json;
  }

  static List<SignalcdListDeploymentResponse> listFromJson(List<dynamic> json) {
    return json == null ? List<SignalcdListDeploymentResponse>() : json.map((value) => SignalcdListDeploymentResponse.fromJson(value)).toList();
  }

  static Map<String, SignalcdListDeploymentResponse> mapFromJson(Map<String, dynamic> json) {
    var map = Map<String, SignalcdListDeploymentResponse>();
    if (json != null && json.isNotEmpty) {
      json.forEach((String key, dynamic value) => map[key] = SignalcdListDeploymentResponse.fromJson(value));
    }
    return map;
  }

  // maps a json object with a list of SignalcdListDeploymentResponse-objects as value to a dart map
  static Map<String, List<SignalcdListDeploymentResponse>> mapListFromJson(Map<String, dynamic> json) {
    var map = Map<String, List<SignalcdListDeploymentResponse>>();
     if (json != null && json.isNotEmpty) {
       json.forEach((String key, dynamic value) {
         map[key] = SignalcdListDeploymentResponse.listFromJson(value);
       });
     }
     return map;
  }
}

