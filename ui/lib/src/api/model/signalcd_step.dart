part of openapi.api;

class SignalcdStep {
  
  String name = null;
  
  String image = null;
  
  List<String> imagePullSecrets = [];
  
  List<String> commands = [];
  
  SignalcdStatus status = null;
  SignalcdStep();

  @override
  String toString() {
    return 'SignalcdStep[name=$name, image=$image, imagePullSecrets=$imagePullSecrets, commands=$commands, status=$status, ]';
  }

  SignalcdStep.fromJson(Map<String, dynamic> json) {
    if (json == null) return;
    name = json['name'];
    image = json['image'];
    imagePullSecrets = (json['ImagePullSecrets'] == null) ?
      null :
      (json['ImagePullSecrets'] as List).cast<String>();
    commands = (json['commands'] == null) ?
      null :
      (json['commands'] as List).cast<String>();
    status = (json['status'] == null) ?
      null :
      SignalcdStatus.fromJson(json['status']);
  }

  Map<String, dynamic> toJson() {
    Map <String, dynamic> json = {};
    if (name != null)
      json['name'] = name;
    if (image != null)
      json['image'] = image;
    if (imagePullSecrets != null)
      json['ImagePullSecrets'] = imagePullSecrets;
    if (commands != null)
      json['commands'] = commands;
    if (status != null)
      json['status'] = status;
    return json;
  }

  static List<SignalcdStep> listFromJson(List<dynamic> json) {
    return json == null ? List<SignalcdStep>() : json.map((value) => SignalcdStep.fromJson(value)).toList();
  }

  static Map<String, SignalcdStep> mapFromJson(Map<String, dynamic> json) {
    var map = Map<String, SignalcdStep>();
    if (json != null && json.isNotEmpty) {
      json.forEach((String key, dynamic value) => map[key] = SignalcdStep.fromJson(value));
    }
    return map;
  }

  // maps a json object with a list of SignalcdStep-objects as value to a dart map
  static Map<String, List<SignalcdStep>> mapListFromJson(Map<String, dynamic> json) {
    var map = Map<String, List<SignalcdStep>>();
     if (json != null && json.isNotEmpty) {
       json.forEach((String key, dynamic value) {
         map[key] = SignalcdStep.listFromJson(value);
       });
     }
     return map;
  }
}

