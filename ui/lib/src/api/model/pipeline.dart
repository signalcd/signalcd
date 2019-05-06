part of swagger.api;

class Pipeline {
  
  String id = null;
  

  Step steps = null;
  

  Check checks = null;
  
  Pipeline();

  @override
  String toString() {
    return 'Pipeline[id=$id, steps=$steps, checks=$checks, ]';
  }

  Pipeline.fromJson(Map<String, dynamic> json) {
    if (json == null) return;
    id =
        json['id']
    ;
    steps =
      
      
      new Step.fromJson(json['steps'])
;
    checks =
      
      
      new Check.fromJson(json['checks'])
;
  }

  Map<String, dynamic> toJson() {
    return {
      'id': id,
      'steps': steps,
      'checks': checks
     };
  }

  static List<Pipeline> listFromJson(List<dynamic> json) {
    return json == null ? new List<Pipeline>() : json.map((value) => new Pipeline.fromJson(value)).toList();
  }

  static Map<String, Pipeline> mapFromJson(Map<String, Map<String, dynamic>> json) {
    var map = new Map<String, Pipeline>();
    if (json != null && json.length > 0) {
      json.forEach((String key, Map<String, dynamic> value) => map[key] = new Pipeline.fromJson(value));
    }
    return map;
  }
}

