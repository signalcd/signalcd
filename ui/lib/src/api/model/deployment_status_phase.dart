part of openapi.api;

class DeploymentStatusPhase {
  /// The underlying value of this enum member.
  final String value;

  const DeploymentStatusPhase._internal(this.value);

  static const DeploymentStatusPhase dONTUSETHISVALUE_ = const DeploymentStatusPhase._internal("DONTUSETHISVALUE");
  static const DeploymentStatusPhase uNKNOWN_ = const DeploymentStatusPhase._internal("UNKNOWN");
  static const DeploymentStatusPhase sUCCESS_ = const DeploymentStatusPhase._internal("SUCCESS");
  static const DeploymentStatusPhase fAILURE_ = const DeploymentStatusPhase._internal("FAILURE");
  static const DeploymentStatusPhase pROGRESS_ = const DeploymentStatusPhase._internal("PROGRESS");
  static const DeploymentStatusPhase pENDING_ = const DeploymentStatusPhase._internal("PENDING");
  static const DeploymentStatusPhase kILLED_ = const DeploymentStatusPhase._internal("KILLED");

  static DeploymentStatusPhase fromJson(String value) {
    return new DeploymentStatusPhaseTypeTransformer().decode(value);
  }
}

class DeploymentStatusPhaseTypeTransformer {

  dynamic encode(DeploymentStatusPhase data) {
    return data.value;
  }

  DeploymentStatusPhase decode(dynamic data) {
    switch (data) {
      case "DONTUSETHISVALUE": return DeploymentStatusPhase.dONTUSETHISVALUE_;
      case "UNKNOWN": return DeploymentStatusPhase.uNKNOWN_;
      case "SUCCESS": return DeploymentStatusPhase.sUCCESS_;
      case "FAILURE": return DeploymentStatusPhase.fAILURE_;
      case "PROGRESS": return DeploymentStatusPhase.pROGRESS_;
      case "PENDING": return DeploymentStatusPhase.pENDING_;
      case "KILLED": return DeploymentStatusPhase.kILLED_;
      default: throw('Unknown enum value to decode: $data');
    }
  }
}

