part of openapi.api;

class SignalcdCheck {
  
  String name = null;
  
  String image = null;
  
  List<String> imagePullSecrets = [];
  
  String duration = null;
  SignalcdCheck();

  @override
  String toString() {
    return 'SignalcdCheck[name=$name, image=$image, imagePullSecrets=$imagePullSecrets, duration=$duration, ]';
  }

  SignalcdCheck.fromJson(Map<String, dynamic> json) {
    if (json == null) return;
    name = json['name'];
    image = json['image'];
    imagePullSecrets = (json['ImagePullSecrets'] == null) ?
      null :
      (json['ImagePullSecrets'] as List).cast<String>();
    duration = json['duration'];
  }

  Map<String, dynamic> toJson() {
    Map <String, dynamic> json = {};
    if (name != null)
      json['name'] = name;
    if (image != null)
      json['image'] = image;
    if (imagePullSecrets != null)
      json['ImagePullSecrets'] = imagePullSecrets;
    if (duration != null)
      json['duration'] = duration;
    return json;
  }

  static List<SignalcdCheck> listFromJson(List<dynamic> json) {
    return json == null ? List<SignalcdCheck>() : json.map((value) => SignalcdCheck.fromJson(value)).toList();
  }

  static Map<String, SignalcdCheck> mapFromJson(Map<String, dynamic> json) {
    var map = Map<String, SignalcdCheck>();
    if (json != null && json.isNotEmpty) {
      json.forEach((String key, dynamic value) => map[key] = SignalcdCheck.fromJson(value));
    }
    return map;
  }

  // maps a json object with a list of SignalcdCheck-objects as value to a dart map
  static Map<String, List<SignalcdCheck>> mapListFromJson(Map<String, dynamic> json) {
    var map = Map<String, List<SignalcdCheck>>();
     if (json != null && json.isNotEmpty) {
       json.forEach((String key, dynamic value) {
         map[key] = SignalcdCheck.listFromJson(value);
       });
     }
     return map;
  }
}

