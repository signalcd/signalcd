part of swagger.api;

class SignalcdDeploymentStatus {
  
  DeploymentStatusPhase phase = null;
  
  SignalcdDeploymentStatus();

  @override
  String toString() {
    return 'SignalcdDeploymentStatus[phase=$phase, ]';
  }

  SignalcdDeploymentStatus.fromJson(Map<String, dynamic> json) {
    if (json == null) return;
    phase =
      
      
      new DeploymentStatusPhase.fromJson(json['phase'])
;
  }

  Map<String, dynamic> toJson() {
    return {
      'phase': phase
     };
  }

  static List<SignalcdDeploymentStatus> listFromJson(List<dynamic> json) {
    return json == null ? new List<SignalcdDeploymentStatus>() : json.map((value) => new SignalcdDeploymentStatus.fromJson(value)).toList();
  }

  static Map<String, SignalcdDeploymentStatus> mapFromJson(Map<String, Map<String, dynamic>> json) {
    var map = new Map<String, SignalcdDeploymentStatus>();
    if (json != null && json.length > 0) {
      json.forEach((String key, Map<String, dynamic> value) => map[key] = new SignalcdDeploymentStatus.fromJson(value));
    }
    return map;
  }
}

