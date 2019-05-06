part of swagger.api;

class Options {
  
  String agent = null;
  

  bool healthy = null;
  
  Options();

  @override
  String toString() {
    return 'Options[agent=$agent, healthy=$healthy, ]';
  }

  Options.fromJson(Map<String, dynamic> json) {
    if (json == null) return;
    agent =
        json['agent']
    ;
    healthy =
        json['healthy']
    ;
  }

  Map<String, dynamic> toJson() {
    return {
      'agent': agent,
      'healthy': healthy
     };
  }

  static List<Options> listFromJson(List<dynamic> json) {
    return json == null ? new List<Options>() : json.map((value) => new Options.fromJson(value)).toList();
  }

  static Map<String, Options> mapFromJson(Map<String, Map<String, dynamic>> json) {
    var map = new Map<String, Options>();
    if (json != null && json.length > 0) {
      json.forEach((String key, Map<String, dynamic> value) => map[key] = new Options.fromJson(value));
    }
    return map;
  }
}

