part of swagger.api;

class Step {
  
  String name = null;
  

  String image = null;
  

  List<String> imagePullSecrets = [];
  

  List<String> commands = [];
  
  Step();

  @override
  String toString() {
    return 'Step[name=$name, image=$image, imagePullSecrets=$imagePullSecrets, commands=$commands, ]';
  }

  Step.fromJson(Map<String, dynamic> json) {
    if (json == null) return;
    name =
        json['name']
    ;
    image =
        json['image']
    ;
    imagePullSecrets =
        (json['imagePullSecrets'] as List).map((item) => item as String).toList()
    ;
    commands =
        (json['commands'] as List).map((item) => item as String).toList()
    ;
  }

  Map<String, dynamic> toJson() {
    return {
      'name': name,
      'image': image,
      'imagePullSecrets': imagePullSecrets,
      'commands': commands
     };
  }

  static List<Step> listFromJson(List<dynamic> json) {
    return json == null ? new List<Step>() : json.map((value) => new Step.fromJson(value)).toList();
  }

  static Map<String, Step> mapFromJson(Map<String, Map<String, dynamic>> json) {
    var map = new Map<String, Step>();
    if (json != null && json.length > 0) {
      json.forEach((String key, Map<String, dynamic> value) => map[key] = new Step.fromJson(value));
    }
    return map;
  }
}

