part of swagger.api;

class Deploymentstatus {
  
  String phase = null;
  //enum phaseEnum {  unknown,  success,  failed,  progress,  };
  Deploymentstatus();

  @override
  String toString() {
    return 'Deploymentstatus[phase=$phase, ]';
  }

  Deploymentstatus.fromJson(Map<String, dynamic> json) {
    if (json == null) return;
    phase =
        json['phase']
    ;
  }

  Map<String, dynamic> toJson() {
    return {
      'phase': phase
     };
  }

  static List<Deploymentstatus> listFromJson(List<dynamic> json) {
    return json == null ? new List<Deploymentstatus>() : json.map((value) => new Deploymentstatus.fromJson(value)).toList();
  }

  static Map<String, Deploymentstatus> mapFromJson(Map<String, Map<String, dynamic>> json) {
    var map = new Map<String, Deploymentstatus>();
    if (json != null && json.length > 0) {
      json.forEach((String key, Map<String, dynamic> value) => map[key] = new Deploymentstatus.fromJson(value));
    }
    return map;
  }
}

