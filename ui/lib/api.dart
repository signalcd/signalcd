import 'package:ui/src/api/api.dart';

class API {
  API() {
    ApiClient client = ApiClient(basePath: "/api");

    this.pipelines = PipelineApi(client);
  }

  PipelineApi pipelines;
}
