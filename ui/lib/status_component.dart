import 'package:angular/angular.dart';

@Component(
  selector: 'status-icon',
  templateUrl: 'status_component.html',
  directives: [coreDirectives],
)
class StatusComponent implements OnInit {
  @Input()
  String status;

  @override
  void ngOnInit() {
    print(status);
  }
}
