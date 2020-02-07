part of openapi.api;

class SignalcdGetCurrentDeploymentResponse {
  
  SignalcdDeployment deployment = null;
  SignalcdGetCurrentDeploymentResponse();

  @override
  String toString() {
    return 'SignalcdGetCurrentDeploymentResponse[deployment=$deployment, ]';
  }

  SignalcdGetCurrentDeploymentResponse.fromJson(Map<String, dynamic> json) {
    if (json == null) return;
    deployment = (json['deployment'] == null) ?
      null :
      SignalcdDeployment.fromJson(json['deployment']);
  }

  Map<String, dynamic> toJson() {
    Map <String, dynamic> json = {};
    if (deployment != null)
      json['deployment'] = deployment;
    return json;
  }

  static List<SignalcdGetCurrentDeploymentResponse> listFromJson(List<dynamic> json) {
    return json == null ? List<SignalcdGetCurrentDeploymentResponse>() : json.map((value) => SignalcdGetCurrentDeploymentResponse.fromJson(value)).toList();
  }

  static Map<String, SignalcdGetCurrentDeploymentResponse> mapFromJson(Map<String, dynamic> json) {
    var map = Map<String, SignalcdGetCurrentDeploymentResponse>();
    if (json != null && json.isNotEmpty) {
      json.forEach((String key, dynamic value) => map[key] = SignalcdGetCurrentDeploymentResponse.fromJson(value));
    }
    return map;
  }

  // maps a json object with a list of SignalcdGetCurrentDeploymentResponse-objects as value to a dart map
  static Map<String, List<SignalcdGetCurrentDeploymentResponse>> mapListFromJson(Map<String, dynamic> json) {
    var map = Map<String, List<SignalcdGetCurrentDeploymentResponse>>();
     if (json != null && json.isNotEmpty) {
       json.forEach((String key, dynamic value) {
         map[key] = SignalcdGetCurrentDeploymentResponse.listFromJson(value);
       });
     }
     return map;
  }
}

