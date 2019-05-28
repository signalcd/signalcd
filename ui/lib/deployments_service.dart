import 'package:angular/angular.dart';
import 'package:ui/api.dart';
import 'package:ui/src/api/api.dart';

@Injectable()
class DeploymentsService {
  final API _api;

  DeploymentsService(this._api);

  Future<List<Deployment>> deployments() async {
    return await _api.deployments.deployments();
  }

  Future<Deployment> deploy(String pipelineID) async {
    return await _api.deployments.setCurrentDeployment(pipelineID);
  }
}
