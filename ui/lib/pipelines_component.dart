import 'package:angular/angular.dart';
import 'package:ui/deployments_service.dart';
import 'package:ui/pipelines_service.dart';
import 'package:ui/src/api/api.dart';
import 'package:ui/timeago_component.dart';

@Component(
  selector: 'pipelines-list',
  templateUrl: 'pipelines_component.html',
  directives: [
    coreDirectives,
    TimeagoComponent,
  ],
  providers: [
    DeploymentsService,
    PipelinesService,
  ],
)
class PipelinesComponent implements OnInit {
  final DeploymentsService _deploymentsService;
  final PipelinesService _pipelinesService;

  PipelinesComponent(this._deploymentsService, this._pipelinesService);

  List<SignalcdPipeline> pipelines = [];

  @override
  void ngOnInit() {
    _pipelinesService
        .pipelines()
        .then((List<SignalcdPipeline> pipelines) => this.pipelines = pipelines);
  }

  void deploy(SignalcdPipeline pipeline) {
    _deploymentsService
        .deploy(pipeline.id)
        .then((SignalcdDeployment deployment) =>
            print('pipeline ${deployment.number} deployed!'))
        .catchError(() => print('error deploying pipeline'));
  }
}
