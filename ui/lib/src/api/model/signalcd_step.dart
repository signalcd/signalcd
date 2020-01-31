part of swagger.api;

class SignalcdStep {
  
  String name = null;
  

  String image = null;
  

  List<String> imagePullSecrets = [];
  

  List<String> commands = [];
  
  SignalcdStep();

  @override
  String toString() {
    return 'SignalcdStep[name=$name, image=$image, imagePullSecrets=$imagePullSecrets, commands=$commands, ]';
  }

  SignalcdStep.fromJson(Map<String, dynamic> json) {
    if (json == null) return;
    name =
        json['name']
    ;
    image =
        json['image']
    ;
    imagePullSecrets =
        (json['ImagePullSecrets'] as List).map((item) => item as String).toList()
    ;
    commands =
        (json['commands'] as List).map((item) => item as String).toList()
    ;
  }

  Map<String, dynamic> toJson() {
    return {
      'name': name,
      'image': image,
      'ImagePullSecrets': imagePullSecrets,
      'commands': commands
     };
  }

  static List<SignalcdStep> listFromJson(List<dynamic> json) {
    return json == null ? new List<SignalcdStep>() : json.map((value) => new SignalcdStep.fromJson(value)).toList();
  }

  static Map<String, SignalcdStep> mapFromJson(Map<String, Map<String, dynamic>> json) {
    var map = new Map<String, SignalcdStep>();
    if (json != null && json.length > 0) {
      json.forEach((String key, Map<String, dynamic> value) => map[key] = new SignalcdStep.fromJson(value));
    }
    return map;
  }
}

