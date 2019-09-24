part of swagger.api;

class Pipeline {
  
  String id = null;
  

  String name = null;
  

  List<Step> steps = [];
  

  List<Check> checks = [];
  

  DateTime created = null;
  
  Pipeline();

  @override
  String toString() {
    return 'Pipeline[id=$id, name=$name, steps=$steps, checks=$checks, created=$created, ]';
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
    created = json['created'] == null ? null : DateTime.parse(json['created']);
  }

  Map<String, dynamic> toJson() {
    return {
      'id': id,
      'name': name,
      'steps': steps,
      'checks': checks,
      'created': created == null ? '' : created.toUtc().toIso8601String()
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

