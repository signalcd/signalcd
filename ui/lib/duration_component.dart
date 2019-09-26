import 'dart:async';

import 'package:angular/angular.dart';

@Component(
  selector: 'signalcd-duration',
  template: '<span>{{ duration }}</span>',
)
class DurationComponent implements OnInit, OnDestroy {
  @Input('start')
  DateTime start;

  @Input('end')
  DateTime end;

  Timer timer;
  String duration = '';

  @override
  void ngOnInit() {
    this.timer = Timer.periodic(const Duration(seconds: 1), tick);
    this.tick(this.timer);
  }

  @override
  void ngOnDestroy() {
    this.timer.cancel();
  }

  void tick(Timer timer) {
    this.duration = this.format(this.start, this.end);
  }

  String format(DateTime start, DateTime end) {
    if (start.isBefore(DateTime(1900))) {
      return '';
    }

    Duration diff = end.difference(start);

    if (end.isBefore(DateTime(1900))) {
      diff = DateTime.now().difference(start);
    }

    String min = diff.inMinutes.toString();
    if (diff.inMinutes < 10) {
      min = '0${diff.inMinutes}';
    }

    String sec = diff.inSeconds.toString();
    if (diff.inSeconds < 10) {
      sec = '0${diff.inSeconds}';
    }

    return '$min:$sec';
  }
}
