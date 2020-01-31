import 'package:angular/angular.dart';
import 'package:ui/api.dart';
import 'package:ui/src/api/api.dart';

@Injectable()
class DeploymentsService {
  final API _api;

  DeploymentsService(this._api);

  Future<List<SignalcdDeployment>> deployments() async {
    return await _api.ui.listDeployment().then((resp) => resp.deployments);
  }

  Future<SignalcdDeployment> deploy(String pipelineID) async {
    return await _api.ui
        .setCurrentDeployment(pipelineID)
        .then((resp) => resp.deployment);
  }
}
