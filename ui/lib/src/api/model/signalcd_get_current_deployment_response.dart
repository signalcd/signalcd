part of swagger.api;

class SignalcdGetCurrentDeploymentResponse {
  
  SignalcdDeployment deployment = null;
  
  SignalcdGetCurrentDeploymentResponse();

  @override
  String toString() {
    return 'SignalcdGetCurrentDeploymentResponse[deployment=$deployment, ]';
  }

  SignalcdGetCurrentDeploymentResponse.fromJson(Map<String, dynamic> json) {
    if (json == null) return;
    deployment =
      
      
      new SignalcdDeployment.fromJson(json['deployment'])
;
  }

  Map<String, dynamic> toJson() {
    return {
      'deployment': deployment
     };
  }

  static List<SignalcdGetCurrentDeploymentResponse> listFromJson(List<dynamic> json) {
    return json == null ? new List<SignalcdGetCurrentDeploymentResponse>() : json.map((value) => new SignalcdGetCurrentDeploymentResponse.fromJson(value)).toList();
  }

  static Map<String, SignalcdGetCurrentDeploymentResponse> mapFromJson(Map<String, Map<String, dynamic>> json) {
    var map = new Map<String, SignalcdGetCurrentDeploymentResponse>();
    if (json != null && json.length > 0) {
      json.forEach((String key, Map<String, dynamic> value) => map[key] = new SignalcdGetCurrentDeploymentResponse.fromJson(value));
    }
    return map;
  }
}

