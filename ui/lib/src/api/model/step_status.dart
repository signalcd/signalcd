part of swagger.api;

class StepStatus {
  
  String logs = null;
  
  StepStatus();

  @override
  String toString() {
    return 'StepStatus[logs=$logs, ]';
  }

  StepStatus.fromJson(Map<String, dynamic> json) {
    if (json == null) return;
    logs =
        json['logs']
    ;
  }

  Map<String, dynamic> toJson() {
    return {
      'logs': logs
     };
  }

  static List<StepStatus> listFromJson(List<dynamic> json) {
    return json == null ? new List<StepStatus>() : json.map((value) => new StepStatus.fromJson(value)).toList();
  }

  static Map<String, StepStatus> mapFromJson(Map<String, Map<String, dynamic>> json) {
    var map = new Map<String, StepStatus>();
    if (json != null && json.length > 0) {
      json.forEach((String key, Map<String, dynamic> value) => map[key] = new StepStatus.fromJson(value));
    }
    return map;
  }
}

