import 'package:ui/src/api/api.dart';

class API {
  API() {
    ApiClient client = ApiClient(basePath: "/api/v1");

    this.deployments = DeploymentsApi(client);
    this.pipelines = PipelineApi(client);
  }

  DeploymentsApi deployments;
  PipelineApi pipelines;
}
