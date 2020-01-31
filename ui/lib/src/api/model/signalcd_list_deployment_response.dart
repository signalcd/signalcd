part of swagger.api;

class SignalcdListDeploymentResponse {
  
  List<SignalcdDeployment> deployments = [];
  
  SignalcdListDeploymentResponse();

  @override
  String toString() {
    return 'SignalcdListDeploymentResponse[deployments=$deployments, ]';
  }

  SignalcdListDeploymentResponse.fromJson(Map<String, dynamic> json) {
    if (json == null) return;
    deployments =
      SignalcdDeployment.listFromJson(json['deployments'])
;
  }

  Map<String, dynamic> toJson() {
    return {
      'deployments': deployments
     };
  }

  static List<SignalcdListDeploymentResponse> listFromJson(List<dynamic> json) {
    return json == null ? new List<SignalcdListDeploymentResponse>() : json.map((value) => new SignalcdListDeploymentResponse.fromJson(value)).toList();
  }

  static Map<String, SignalcdListDeploymentResponse> mapFromJson(Map<String, Map<String, dynamic>> json) {
    var map = new Map<String, SignalcdListDeploymentResponse>();
    if (json != null && json.length > 0) {
      json.forEach((String key, Map<String, dynamic> value) => map[key] = new SignalcdListDeploymentResponse.fromJson(value));
    }
    return map;
  }
}

