import 'dart:async';

import 'package:angular/angular.dart';
import 'package:ui/deployments_service.dart';
import 'package:ui/src/api/api.dart';
import 'package:ui/status_component.dart';
import 'package:ui/timeago_pipe.dart';

@Component(
  selector: 'deployments-list',
  templateUrl: 'deployments_component.html',
  directives: [
    coreDirectives,
    StatusComponent,
  ],
  providers: [
    DeploymentsService,
  ],
  pipes: [
    TimeagoPipe,
  ],
)
class DeploymentsComponent implements OnInit, OnDestroy {
  final DeploymentsService _deploymentsService;

  DeploymentsComponent(this._deploymentsService);

  Timer timer;
  List<Deployment> deployments = [];

  @override
  void ngOnInit() {
    getDeployments();

    timer = Timer.periodic(
      Duration(seconds: 1),
      (Timer timer) => getDeployments(),
    );
  }

  @override
  void ngOnDestroy() {
    timer.cancel();
  }

  void getDeployments() {
    _deploymentsService.deployments().then((List<Deployment> deployments) {
      // Only update if number of deployments changed
      if (this.deployments.length != deployments.length) {
        this.deployments = deployments;
      }
    });
  }
}
