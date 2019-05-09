part of swagger.api;

class Pipeline {
  
  String id = null;
  

  String name = null;
  

  List<Step> steps = [];
  

  List<Check> checks = [];
  
  Pipeline();

  @override
  String toString() {
    return 'Pipeline[id=$id, name=$name, steps=$steps, checks=$checks, ]';
  }

  Pipeline.fromJson(Map<String, dynamic> json) {
    if (json == null) return;
    id =
        json['id']
    ;
    name =
        json['name']
    ;
    steps =
      Step.listFromJson(json['steps'])
;
    checks =
      Check.listFromJson(json['checks'])
;
  }

  Map<String, dynamic> toJson() {
    return {
      'id': id,
      'name': name,
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

