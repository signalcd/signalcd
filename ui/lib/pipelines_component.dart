import 'package:angular/angular.dart';
import 'package:ui/pipelines_service.dart';
import 'package:ui/src/api/api.dart';

@Component(
  selector: 'pipelines-list',
  templateUrl: 'pipelines_component.html',
  directives: [
    coreDirectives,
  ],
  providers: [
    PipelinesService,
  ],
)
class PipelinesComponent implements OnInit {
  final PipelinesService _pipelinesService;

  PipelinesComponent(this._pipelinesService);

  List<Pipeline> pipelines = [];

  @override
  void ngOnInit() {
    _pipelinesService
        .pipelines()
        .then((List<Pipeline> pipelines) => this.pipelines = pipelines);
  }

  void deploy(Pipeline pipeline) {
    _pipelinesService
        .deploy(pipeline.id)
        .then((dynamic) => print('pipeline ${pipeline.id} deployed!'))
        .catchError(() => print('error dpeloying pipeline'));
  }
}
